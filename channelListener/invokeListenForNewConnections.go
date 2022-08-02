package channelListener

import (
	"go.uber.org/fx"
	"golang.org/x/net/context"
)

func invokeListenForNewConnections() fx.Option {
	return fx.Invoke(
		func(
			params struct {
				fx.In
				NetManager *manager
				CancelFunc context.CancelFunc
				Lifecycle  fx.Lifecycle
			},
		) {
			params.Lifecycle.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return params.NetManager.ListenForNewConnections()
				},
				OnStop: func(ctx context.Context) error {
					params.CancelFunc()
					return nil
				},
			})
		},
	)
}
