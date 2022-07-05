package internal

import (
	"github.com/bhbosman/gocomms/intf"
	"go.uber.org/fx"
	"golang.org/x/crypto/ssh"
)

type ISshConnectionReactorChannelManager interface {
	AddAcceptedChannel(uniqueReference string, acceptedChannel *fx.App) error
	RemoveAcceptedChannel(uniqueReference string) error
}

type ISshConnectionReactor interface {
	intf.IConnectionReactor
	ISshConnectionReactorChannelManager
	CanAcceptChannel(name string) (bool, ssh.RejectionReason, string, error)
}
