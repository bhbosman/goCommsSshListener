package stack

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goCommsSshListener/common"
	"github.com/bhbosman/goCommsSshListener/internal"
	"github.com/bhbosman/goConnectionManager"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/services/interfaces"
	common2 "github.com/bhbosman/gocomms/common"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/zap"
	"net"
	"reflect"
)

func createStackState(
	connectionType model.ConnectionType,
	ConnectionManager goConnectionManager.IService,
	UniqueSessionNumber interfaces.IUniqueReferenceService,
	connectionReactor internal.ISshConnectionReactor,
	logger *zap.Logger,
	conn net.Conn,
	ctx context.Context,
	ctxCancelFunc context.CancelFunc,
	SshChannelSettings common.ISshChannelSettings,
	goFunctionCounter GoFunctionCounter.IService,

) common2.IStackState {
	return common2.NewStackState(
		goCommsDefinitions.SshStackName,
		true,
		func() (common2.IStackCreateData, error) {
			return NewStackData(
				connectionType,
				conn,
				ctx,
				ctxCancelFunc,
				logger,
				ConnectionManager,
				UniqueSessionNumber,
				connectionReactor,
				SshChannelSettings,
				goFunctionCounter,
			), nil
		},
		func(
			connectionType model.ConnectionType,
			stackData common2.IStackCreateData,
		) error {
			if closer, ok := stackData.(*data); ok {
				return closer.Close()
			}
			return WrongStackDataError(connectionType, stackData)
		},
		func(
			inputStreamForStack common2.IInputStreamForStack,
			stackData common2.IStackCreateData,
			ToReactorFunc rxgo.NextFunc,
		) (common2.IInputStreamForStack, error) {
			if stackDataInstance, ok := stackData.(*data); ok {
				return stackDataInstance.Start(ctx, ToReactorFunc)
			}
			return nil, WrongStackDataError(connectionType, stackData)
		},
		func(
			stackData interface{},
		) error {
			if stop, ok := stackData.(*data); ok {
				return stop.Stop()
			}
			return WrongStackDataError(connectionType, stackData)
		},
	)
}

func WrongStackDataError(connectionType model.ConnectionType, stackData interface{}) error {
	return common2.NewWrongStackDataType(
		goCommsDefinitions.SshStackName,
		connectionType,
		reflect.TypeOf((*data)(nil)),
		reflect.TypeOf(stackData))
}
