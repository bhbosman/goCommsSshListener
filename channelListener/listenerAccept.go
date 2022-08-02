package channelListener

import (
	"context"
	"github.com/bhbosman/goCommsSshListener/common"
)

type listenerAccept struct {
	listener interface {
		accept() (common.INewChannel, error)
	}
}

func (self *listenerAccept) acceptWithContext() (common.INewChannel, context.CancelFunc, error) {
	conn, err := self.listener.accept()
	if err != nil {
		return nil, nil, err
	}
	return conn,
		func() {
			// do nothing
		},
		nil
}
