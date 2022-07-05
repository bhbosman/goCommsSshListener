package channelListener

import (
	"go.uber.org/fx"
	"net"
)

type sshChannelListenerOverrideListenerFactory struct {
	listenerFactory func() (net.Listener, error)
}

func (self *sshChannelListenerOverrideListenerFactory) apply(settings *sshChannelListenerManagerSettings) (fx.Option, error) {
	settings.setListenerFactory(self.listenerFactory)
	return nil, nil
}

func newSshOverrideListener(listenerFactory func() (net.Listener, error)) *sshChannelListenerOverrideListenerFactory {
	return &sshChannelListenerOverrideListenerFactory{listenerFactory: listenerFactory}
}
