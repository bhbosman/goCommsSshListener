package channelListener

import (
	"go.uber.org/fx"
	"golang.org/x/crypto/ssh"
)

func provideNewChannelChannel(channels <-chan ssh.NewChannel) fx.Option {
	return fx.Provide(
		fx.Annotated{
			Target: func() (<-chan ssh.NewChannel, error) {
				return channels, nil
			},
		},
	)
}
