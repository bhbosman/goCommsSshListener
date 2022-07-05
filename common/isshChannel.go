package common

import (
	"github.com/bhbosman/goCommsDefinitions"
	"golang.org/x/crypto/ssh"
)

type ISshChannel interface {
	ssh.Channel
	goCommsDefinitions.ISpecificInformationForConnection
	ChannelType() string
}
