package channelListener

import (
	"github.com/bhbosman/gocomms/common"
	"net"
)

type sshChannelListenerOverrideListenerAcceptFactory struct {
	listenerAcceptFactory func(listener net.Listener) (ISshListenerAccept, error)
}

func (self *sshChannelListenerOverrideListenerAcceptFactory) ApplyNetManagerSettings(settings *common.NetManagerSettings) error {
	return nil
}

func (self *sshChannelListenerOverrideListenerAcceptFactory) apply(settings *sshChannelListenerManagerSettings) error {
	settings.setListenerAcceptFactory(self.listenerAcceptFactory)
	return nil
}

func newSshOverrideListenerAcceptFactory(listenerAcceptFactory func(listener net.Listener) (ISshListenerAccept, error)) *sshChannelListenerOverrideListenerAcceptFactory {
	return &sshChannelListenerOverrideListenerAcceptFactory{
		listenerAcceptFactory: listenerAcceptFactory,
	}
}
