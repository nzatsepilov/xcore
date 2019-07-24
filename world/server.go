package world

import (
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
	config    *config.Config
	db        *db.DB
	tcpServer net.TCPServer
}

func NewServer(c *config.Config) (net.Server, error) {
	xdb, err := db.New(c.DBConfig)
	if err != nil {
		return nil, err
	}

	s := new(server)
	s.config = c
	s.db = xdb
	s.tcpServer = net.NewTCPServer(&net.ServerParameters{
		OnConnection: s.handleConnection,
		OnError: func(err error) {
			println(err)
		},
	})

	return s, nil
}

func (srv *server) start() error {
	if err := srv.tcpServer.Start(srv.config.WorldServerAddress); err != nil {
		return err
	}

	log.Printf("World server started, listening `%s`\n", srv.config.WorldServerAddress)
	return nil
}

func (srv *server) stop() error {
	if err := srv.tcpServer.Stop(); err != nil {
		return err
	}

	log.Println("World server stopped")
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
	s := NewSession(conn)
	go s.start()
}

func (srv *server) waitExitSignal() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGINT)
	<-c
}
