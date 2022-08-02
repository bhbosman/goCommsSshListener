package channelListener

import (
	"github.com/bhbosman/gocomms/common"
	"net"
)

type sshChannelListenerOverrideListenerAcceptFactory struct {
	listenerAcceptFactory func(listener net.Listener) (iListenerAccept, error)
}

func (self *sshChannelListenerOverrideListenerAcceptFactory) ApplyNetManagerSettings(settings *common.NetManagerSettings) error {
	return nil
}

func (self *sshChannelListenerOverrideListenerAcceptFactory) apply(settings *channelListenerManagerSettings) error {
	settings.setListenerAcceptFactory(self.listenerAcceptFactory)
	return nil
}

func newSshOverrideListenerAcceptFactory(listenerAcceptFactory func(listener net.Listener) (iListenerAccept, error)) *sshChannelListenerOverrideListenerAcceptFactory {
	return &sshChannelListenerOverrideListenerAcceptFactory{
		listenerAcceptFactory: listenerAcceptFactory,
	}
}
