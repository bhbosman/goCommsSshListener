package channelListener

import (
	"go.uber.org/fx"
	"golang.org/x/crypto/ssh"
)

func provideAcceptedChannelRequestChannel(channel <-chan *ssh.Request) fx.Option {
	return fx.Provide(
		fx.Annotated{
			Target: func() <-chan *ssh.Request {
				return channel
			},
		},
	)
}
