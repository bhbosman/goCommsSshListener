package common

import (
	"github.com/bhbosman/goCommsDefinitions"
	"golang.org/x/crypto/ssh"
)

type INewChannel interface {
	goCommsDefinitions.ISpecificInformationForConnection
	Accept(channelType string, additionalData []byte) (IChannel, <-chan *ssh.Request, error)
	Reject(reason ssh.RejectionReason, message string) error
	ChannelType() string
	ExtraData() []byte
}
