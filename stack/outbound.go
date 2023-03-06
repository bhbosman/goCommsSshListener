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

func Outbound(
	connectionType model.ConnectionType,
	ConnectionCancelFunc model.ConnectionCancelFunc,
	logger *zap.Logger,
	ctx context.Context,
	goFunctionCounter GoFunctionCounter.IService,
	opts ...rxgo.Option,
) common2.BoundResult {
	return func() (common2.IStackBoundDefinition, error) {
		return common2.NewStackBoundDefinition(
				func(
					sd common2.IStackCreateData,
					pipeData common2.IPipeCreateData,
					obs gocommon.IObservable,
				) (gocommon.IObservable, error) {
					if pipeData != nil {
						return nil, goerrors.InvalidParam
					}
					if stackData, ok := sd.(*data); ok {
						nextOutboundChannel := make(chan rxgo.Item)
						var err error
						stackData.outBoundHandler, err = RxHandlers.All2(
							goCommsDefinitions.SshStackName,
							model.StreamDirectionUnknown,
							nextOutboundChannel,
							logger,
							ctx,
							true,
						)
						if err != nil {
							return nil, err
						}

						connWrapper, err := common2.NewConnWrapper(
							stackData.conn,
							stackData.ctx,
							stackData.pipeRead,
							stackData.outBoundHandler.OnSendData,
						)
						if err != nil {
							return nil, err
						}

						err = stackData.setConnWrapper(connWrapper)
						if err != nil {
							return nil, err
						}

						outboundStackHandler, err := NewOutboundStackHandler(stackData)
						if err != nil {
							return nil, err
						}

						nextHandler, err := RxHandlers.NewRxNextHandler2(
							goCommsDefinitions.SshStackName,
							ConnectionCancelFunc,
							outboundStackHandler,
							stackData.outBoundHandler,
							logger)
						if err != nil {
							return nil, err
						}

						rxOverride.ForEach2(
							goCommsDefinitions.SshStackName,
							model.StreamDirectionOutbound,
							obs,
							ctx,
							goFunctionCounter,
							nextHandler,
							opts...,
						)
						nextObs := rxgo.FromChannel(nextOutboundChannel, opts...)
						return nextObs, nil
					}
					return nil, WrongStackDataError(connectionType, sd)
				},
				nil),
			nil
	}
}
