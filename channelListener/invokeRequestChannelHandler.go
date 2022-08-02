package channelListener

import (
	"context"
	"github.com/bhbosman/goCommsSshListener/internal"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
	"golang.org/x/crypto/ssh"
)

func invokeRequestChannelHandler() fx.Option {
	return fx.Invoke(
		func(params struct {
			fx.In
			CancelContext     context.Context
			CancelFunc        context.CancelFunc
			RequestChannel    <-chan *ssh.Request
			Lifecycle         fx.Lifecycle
			ToReactorFunc     rxgo.NextFunc `name:"ForReactor"`
			GoFunctionCounter GoFunctionCounter.IService
		}) error {
			params.Lifecycle.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return internal.GoRequestChannelHandler(
						"InvokeRequestChannelHandler",
						params.CancelContext,
						params.RequestChannel,
						params.ToReactorFunc,
						params.GoFunctionCounter,
					)
				},
				OnStop: func(ctx context.Context) error {
					params.CancelFunc()
					return nil
				},
			})
			return nil
		},
	)
}
