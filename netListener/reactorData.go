package netListener

import (
	"github.com/bhbosman/gocommon/messageRouter"
	"github.com/bhbosman/gocommon/messages"
	"golang.org/x/crypto/ssh"
)

type Data struct {
	MessageRouter messageRouter.IMessageRouter
}

func (self *Data) ShutDown() error {
	return nil
}

func (self *Data) Send(message interface{}) error {
	self.MessageRouter.Route(message)
	return nil
}

func (self *Data) AddAcceptedChannel(uniqueReference string, acceptedChannel messages.IApp) error {
	return nil
}

func (self *Data) RemoveAcceptedChannel(uniqueReference string) error {
	return nil
}

func (self *Data) HandleGlobalRequest(name string, req *ssh.Request) error {
	return nil
}

func (self *Data) CanAcceptChannel(name string) (bool, ssh.RejectionReason, string, error) {
	return true, 0, "", nil
}

func NewData() (IConnectionReactorMessageQueueData, error) {
	result := &Data{
		MessageRouter: messageRouter.NewMessageRouter(),
	}
	result.MessageRouter.Add(result.handleEmptyQueue)
	result.MessageRouter.Add(result.handleSshRequest)
	return result, nil
}

func (self *Data) handleEmptyQueue(msg *messages.EmptyQueue) {
}

func (self *Data) handleSshRequest(msg *ssh.Request) {
	if msg.WantReply {
		msg.Reply(true, nil)
	}
}
