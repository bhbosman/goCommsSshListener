package stack

import (
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/gocommon"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/rxOverride"
	"github.com/bhbosman/gocomms/RxHandlers"
	common2 "github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/goerrors"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

func Inbound(
	connectionType model.ConnectionType,
	ConnectionCancelFunc model.ConnectionCancelFunc,
	logger *zap.Logger,
	ctx context.Context,
	goFunctionCounter GoFunctionCounter.IService,
	opts ...rxgo.Option,
) common2.BoundResult {
	return func() (common2.IStackBoundFactory, error) {
		return common2.NewStackBoundDefinition(
				func(
					stackData common2.IStackCreateData,
					pipeData common2.IPipeCreateData,
					obs gocommon.IObservable,
				) (gocommon.IObservable, error) {
					if pipeData != nil {
						return nil, goerrors.InvalidParam
					}
					if sd, ok := stackData.(*data); ok {
						NextInBoundChannel := make(chan rxgo.Item)
						var err error
						sd.inBoundHandler, err = RxHandlers.All2(
							goCommsDefinitions.SshStackName,
							model.StreamDirectionUnknown,
							NextInBoundChannel,
							logger,
							ctx,
							true,
						)
						if err != nil {
							return nil, err
						}

						inboundStackHandler, err := NewInboundStackHandler(sd)
						if err != nil {
							return nil, err
						}

						nextHandler, err := RxHandlers.NewRxNextHandler2(
							goCommsDefinitions.SshStackName,
							ConnectionCancelFunc,
							inboundStackHandler,
							sd.inBoundHandler,
							logger)
						if err != nil {
							return nil, err
						}

						rxOverride.ForEach2(
							goCommsDefinitions.SshStackName,
							model.StreamDirectionInbound,
							obs,
							ctx,
							goFunctionCounter,
							nextHandler,
							opts...)
						nextObs := rxgo.FromChannel(NextInBoundChannel, opts...)
						return nextObs, nil
					}
					return nil, WrongStackDataError(connectionType, stackData)
				},
				nil),
			nil
	}
}
