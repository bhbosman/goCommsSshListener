package channelListener

import "github.com/bhbosman/gocomms/common"

type iListenAppSettingsApply interface {
	common.INetManagerSettingsApply
	apply(settings *channelListenerManagerSettings) error
}
