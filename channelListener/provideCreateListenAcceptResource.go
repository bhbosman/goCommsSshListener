package channelListener

import (
	"go.uber.org/fx"
)

func provideCreateListenAcceptResource(
	params struct {
		fx.In
		Listener iChannelListener
	}) (ISshListenerAccept, error) {
	return &SshListenerAccept{
		Listener: params.Listener,
	}, nil
}
