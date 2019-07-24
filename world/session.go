package world

import (
	"log"
	xnet "net"
	"xcore/core/net"
)

type Session interface {
	start()
}

type session struct {
	sock *net.Socket
}

func NewSession(c *xnet.TCPConn) Session {
	sock := net.NewSocket(c)
	return &session{
		sock: sock,
	}
}

func (s *session) start() {
	if err := s.authorize(); err != nil {
		log.Printf("can not authorize world session: %v", err)
	}
}

func (s *session) authorize() error {
	if err := s.sock.ReceiveData(); err != nil {
		return err
	}

	build := s.sock.MustReadUInt32()
	loginServerId := s.sock.MustReadUInt32()
	accountName := s.sock.MustReadString()
	challenge := s.sock.MustReadUInt32()
	authHash := s.sock.MustReadBytes(20)

	log.Println(build)
	log.Println(loginServerId)
	log.Println(accountName)
	log.Println(challenge)
	log.Println(authHash)

	return nil
}
