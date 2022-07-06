package channelListener

import (
	"github.com/bhbosman/gocomms/common"
)

type sshChannelListenerUserContextValue struct {
	userContext interface{}
}

func (self *sshChannelListenerUserContextValue) ApplyNetManagerSettings(settings *common.NetManagerSettings) error {
	return nil
}

func (self *sshChannelListenerUserContextValue) apply(settings *channelListenerManagerSettings) error {
	settings.userContext = self.userContext
	return nil
}

func newSshUserContextValue(userContext interface{}) *sshChannelListenerUserContextValue {
	return &sshChannelListenerUserContextValue{userContext: userContext}
}
