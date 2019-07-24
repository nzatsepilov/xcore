package net

import (
	"go.uber.org/atomic"
	"net"
)

type tcpListenerCallbacks struct {
	onConnection func(conn *net.TCPConn)
	onError      func(err error)
}

type tcpListener interface {
	listen(address string) error
	stop() error
}

type asyncTCPListener struct {
	callbacks   *tcpListenerCallbacks
	tcpListener *net.TCPListener

	isClosing atomic.Bool
	closing   chan bool
}

func newTCPListener(callbacks *tcpListenerCallbacks) tcpListener {
	listener := new(asyncTCPListener)
	listener.callbacks = callbacks
	listener.closing = make(chan bool, 1)
	return listener
}

func (listener *asyncTCPListener) listen(address string) error {
	tcpAddress, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		return err
	}

	tcpListener, err := net.ListenTCP("tcp4", tcpAddress)
	if err != nil {
		return err
	}

	listener.tcpListener = tcpListener
	go listener.startListenLoop()
	return nil
}

func (listener *asyncTCPListener) stop() error {
	return listener.stopGracefully()
}

func (listener *asyncTCPListener) stopGracefully() error {
	listener.closing <- true
	if err := listener.tcpListener.Close(); err != nil {
		return err
	}
	<-listener.closing
	listener.tcpListener = nil
	return nil
}

func (listener *asyncTCPListener) startListenLoop() {
loop:
	for {
		conn, err := listener.tcpListener.AcceptTCP()

		select {
		case <-listener.closing:
			break loop
		default:
			if err != nil {
				listener.callbacks.onError(err)
			} else {
				listener.callbacks.onConnection(conn)
			}
		}
	}

	listener.closing <- true
}
