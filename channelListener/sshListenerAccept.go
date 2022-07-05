package channelListener

import (
	"context"
	"github.com/bhbosman/goCommsSshListener/common"
)

type SshListenerAccept struct {
	Listener interface {
		Accept() (common.ISshNewChannel, error)
	}
}

func (self *SshListenerAccept) AcceptWithContext() (common.ISshNewChannel, context.CancelFunc, error) {
	conn, err := self.Listener.Accept()
	if err != nil {
		return nil, nil, err
	}
	return conn,
		func() {
			// do nothing
		},
		nil
}
