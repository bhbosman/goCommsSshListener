package common

import (
	"github.com/bhbosman/goCommsDefinitions"
	"golang.org/x/crypto/ssh"
)

type IChannel interface {
	ssh.Channel
	goCommsDefinitions.ISpecificInformationForConnection
	ChannelType() string
}
