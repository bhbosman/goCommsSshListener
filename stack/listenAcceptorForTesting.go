package stack

import (
	"context"
	"net"
)

type ListenAcceptorForTesting struct {
	listener net.Listener
	onAccept func(conn net.Conn, cancelContext context.Context, cancelFunc context.CancelFunc)
}

func (self *ListenAcceptorForTesting) Accept() (net.Conn, error) {
	return self.listener.Accept()
}

func (self *ListenAcceptorForTesting) Addr() net.Addr {
	return self.listener.Addr()
}

func NewListenAcceptor(listener net.Listener) (*ListenAcceptorForTesting, error) {
	return &ListenAcceptorForTesting{listener: listener}, nil
}

func (self *ListenAcceptorForTesting) AcceptWithContext() (net.Conn, context.CancelFunc, error) {
	accept, err := self.Accept()
	if err != nil {
		return nil, nil, err
	}
	cancelContext, cancelFunc := context.WithCancel(context.Background())
	if self.onAccept != nil {
		self.onAccept(accept, cancelContext, cancelFunc)
	}
	return accept, cancelFunc, nil
}
