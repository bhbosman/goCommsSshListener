package internal

import (
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocomms/intf"
	"golang.org/x/crypto/ssh"
)

type ISshConnectionReactorChannelManager interface {
	AddAcceptedChannel(uniqueReference string, acceptedChannel messages.IApp) error
	RemoveAcceptedChannel(uniqueReference string) error
}

type ISshConnectionReactor interface {
	intf.IConnectionReactor
	ISshConnectionReactorChannelManager
	CanAcceptChannel(name string) (bool, ssh.RejectionReason, string, error)
}
