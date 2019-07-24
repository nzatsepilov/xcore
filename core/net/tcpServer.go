package net

import "net"

type TCPServer interface {
	Start(address string) error
	Stop() error
}
type ServerParameters struct {
	OnConnection func(conn *net.TCPConn)
	OnError      func(err error)
}

type asyncTCPServer struct {
	listener tcpListener
}

func NewTCPServer(p *ServerParameters) TCPServer {
	s := new(asyncTCPServer)
	s.listener = newTCPListener(&tcpListenerCallbacks{
		onConnection: p.OnConnection,
		onError:      p.OnError,
	})
	return s
}

func (s *asyncTCPServer) Start(address string) error {
	return s.listener.listen(address)
}

func (s *asyncTCPServer) Stop() error {
	return s.listener.stop()
}
