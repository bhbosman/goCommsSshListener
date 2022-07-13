package netListener

import (
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/intf"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

func ProvideConnectionReactor() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Target: func(
				params struct {
					fx.In
					CancelCtx            context.Context
					CancelFunc           context.CancelFunc
					ConnectionCancelFunc model.ConnectionCancelFunc
					Logger               *zap.Logger
					//ClientContext        interface{} `name:"UserContext"`
					GoFunctionCounter GoFunctionCounter.IService
				},
			) (intf.IConnectionReactor, error) {
				return NewSshConnectionReactor(
					params.CancelCtx,
					params.CancelFunc,
					params.ConnectionCancelFunc,
					params.Logger,
					//params.ClientContext,
					params.GoFunctionCounter,
				)
			},
		},
	)
}
