package channelListener

import (
	"github.com/bhbosman/gocomms/common"
	"net"
)

type sshChannelListenerManagerSettings struct {
	common.NetManagerSettings
	userContext           interface{}
	netListenerFactory    interface{} //func() (net.Listener, error)
	listenerAcceptFactory interface{} //func(IListenerAccept, err)
}

func (self *sshChannelListenerManagerSettings) setListenerAcceptFactory(listenerAcceptFactory func(listener net.Listener) (ISshListenerAccept, error)) {
	self.listenerAcceptFactory = listenerAcceptFactory
}

func (self *sshChannelListenerManagerSettings) setListenerFactory(netListenerFactory func() (net.Listener, error)) {
	self.netListenerFactory = netListenerFactory
}
