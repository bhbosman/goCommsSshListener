package channelListener

import "github.com/bhbosman/gocomms/common"

type ISshListenAppSettingsApply interface {
	common.INetManagerSettingsApply
	apply(settings *sshChannelListenerManagerSettings) error
}
