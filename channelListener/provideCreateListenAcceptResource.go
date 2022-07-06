package channelListener

import (
	"go.uber.org/fx"
)

func provideCreateListenAcceptResource(
	params struct {
		fx.In
		Listener iChannelListener
	}) (iListenerAccept, error) {
	return &listenerAccept{
		listener: params.Listener,
	}, nil
}
