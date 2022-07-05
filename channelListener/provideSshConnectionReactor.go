package channelListener

import (
	"github.com/bhbosman/goCommsSshListener/common"
	"github.com/bhbosman/goCommsSshListener/internal"
	"go.uber.org/fx"
)

func provideSshConnectionReactor(sshConnectionReactor internal.ISshConnectionReactor) fx.Option {
	return fx.Provide(
		fx.Annotated{
			Target: func() internal.ISshConnectionReactor {
				return sshConnectionReactor
			},
		},
	)
}

func provideSshCreateChannelProcess(SshChannelSettings common.ISshChannelSettings) fx.Option {
	return fx.Provide(
		fx.Annotated{
			Target: func() common.ISshChannelSettings {
				return SshChannelSettings
			},
		},
	)
}
