package channelListener

import (
	"context"
	"go.uber.org/fx"
)

func sshInvokeListenForNewConnections(
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
}
