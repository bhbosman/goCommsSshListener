package stack

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goCommsSshListener/common"
	"github.com/bhbosman/goCommsSshListener/internal"
	"github.com/bhbosman/goConnectionManager"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/Services/interfaces"
	"github.com/bhbosman/gocommon/model"
	internalComms "github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/goerrors"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"net"
)

func Provide() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Group: "StackDefinition",
			Target: func(
				params struct {
					fx.In
					ConnectionCancelFunc model.ConnectionCancelFunc
					Opts                 []rxgo.Option
					Logger               *zap.Logger
					ConnectionType       model.ConnectionType
					ConnectionManager    goConnectionManager.IService
					UniqueSessionNumber  interfaces.IUniqueReferenceService
					ConnectionReactor    internal.ISshConnectionReactor
					Conn                 net.Conn
					Ctx                  context.Context
					CtxCancelFunc        context.CancelFunc
					SshChannelSettings   common.ISshChannelSettings `optional:"true"`
					GoFunctionCounter    GoFunctionCounter.IService
				},
			) (*internalComms.StackDefinition, error) {
				var errList error = nil
				if params.ConnectionCancelFunc == nil {
					errList = multierr.Append(errList, goerrors.InvalidParam)
				}

				if params.Logger == nil {
					errList = multierr.Append(errList, goerrors.InvalidParam)
				}

				if errList != nil {
					return nil, errList
				}

				return &internalComms.StackDefinition{
					Name: goCommsDefinitions.SshStackName,
					Inbound: internalComms.NewBoundResultImpl(
						Inbound(
							params.ConnectionType,
							params.ConnectionCancelFunc,
							params.Logger,
							params.Ctx,
							params.GoFunctionCounter,
							params.Opts...,
						),
					),
					Outbound: internalComms.NewBoundResultImpl(
						Outbound(
							params.ConnectionType,
							params.ConnectionCancelFunc,
							params.Logger,
							params.Ctx,
							params.GoFunctionCounter,
							params.Opts...,
						),
					),
					StackState: createStackState(
						params.ConnectionType,
						params.ConnectionManager,
						params.UniqueSessionNumber,
						params.ConnectionReactor,
						params.Logger,
						params.Conn,
						params.Ctx,
						params.CtxCancelFunc,
						params.SshChannelSettings,
						params.GoFunctionCounter,
					),
				}, nil
			},
		},
	)
}
