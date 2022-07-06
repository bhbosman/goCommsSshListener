package channelListener

import (
	"github.com/bhbosman/gocomms/common"
	"net"
)

type channelListenerManagerSettings struct {
	common.NetManagerSettings
	userContext           interface{}
	netListenerFactory    interface{} //func() (net.Listener, error)
	listenerAcceptFactory interface{} //func(IListenerAccept, err)
}

func (self *channelListenerManagerSettings) setListenerAcceptFactory(listenerAcceptFactory func(listener net.Listener) (iListenerAccept, error)) {
	self.listenerAcceptFactory = listenerAcceptFactory
}

func (self *channelListenerManagerSettings) setListenerFactory(netListenerFactory func() (net.Listener, error)) {
	self.netListenerFactory = netListenerFactory
}
