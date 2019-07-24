package auth

import (
	uuid "github.com/satori/go.uuid"
	"log"
	xnet "net"
	"os"
	"os/signal"
	"syscall"
	"xcore/config"
	"xcore/core/db"
	"xcore/core/net"
)

type server struct {
	config *config.Config

	db      *db.DB
	accRepo AccountRepository

	tcpServer net.TCPServer
	realmList *realmProvider
}

func NewServer(c *config.Config) (net.Server, error) {
	xdb, err := db.New(c.DBConfig)
	if err != nil {
		return nil, err
	}

	accRepo, err := NewAccountRepository(c, xdb)
	if err != nil {
		return nil, err
	}

	s := new(server)
	s.config = c
	s.db = xdb
	s.accRepo = accRepo
	s.tcpServer = net.NewTCPServer(&net.ServerParameters{
		OnConnection: s.handleConnection,
		OnError:      s.handleError,
	})
	s.realmList = NewRealmProvider(c)

	return s, nil
}

func (srv *server) start() error {
	if err := srv.tcpServer.Start(srv.config.AuthServerAddress); err != nil {
		return err
	}

	log.Printf("auth server started at `%v`\n", srv.config.AuthServerAddress)

	log.Printf("added %v realm(s):", len(srv.config.Realms))
	for _, r := range srv.config.Realms {
		log.Printf("#%v \"%v\" at %v", r.ID, r.Name, r.Address)
	}

	return nil
}

func (srv *server) stop() error {
	if err := srv.tcpServer.Stop(); err != nil {
		return err
	}

	log.Println("auth server stopped")
	return nil
}

func (srv *server) Run() error {
	if err := srv.start(); err != nil {
		return err
	}
	srv.waitExitSignal()
	return srv.stop()
}

func (srv *server) handleConnection(conn *xnet.TCPConn) {
	id := uuid.NewV4().String()
	s := newSession(id, conn, srv.accRepo, srv.realmList)
	go s.authorize()
}

func (srv *server) handleError(err error) {
	log.Println(err)
}

func (srv *server) waitExitSignal() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGINT)
	<-c
}
