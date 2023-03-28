package netListener

import (
	"github.com/bhbosman/goCommsSshListener/internal"
	"github.com/bhbosman/gocommon/services/IDataShutDown"
	"github.com/bhbosman/gocommon/services/IFxService"
	"github.com/bhbosman/gocommon/services/ISendMessage"
	"golang.org/x/crypto/ssh"
)

type IConnectionReactorMessageQueue interface {
	ISendMessage.ISendMessage
	internal.ISshConnectionReactorChannelManager
	HandleGlobalRequest(name string, req *ssh.Request) error
	CanAcceptChannel(name string) (bool, ssh.RejectionReason, string, error)
}

type IConnectionReactorMessageQueueService interface {
	IConnectionReactorMessageQueue
	IFxService.IFxServices
}

type IConnectionReactorMessageQueueData interface {
	IConnectionReactorMessageQueue
	IDataShutDown.IDataShutDown
}
