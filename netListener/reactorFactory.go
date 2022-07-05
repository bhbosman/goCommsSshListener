package netListener

import (
	"context"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/intf"
	"go.uber.org/zap"
)

type factory struct {
	goFunctionCounter GoFunctionCounter.IService
}

func (self *factory) Create(
	cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	connectionCancelFunc model.ConnectionCancelFunc,
	logger *zap.Logger,
	userContext interface{},
) (intf.IConnectionReactor, error) {
	return NewSshConnectionReactor(
		cancelCtx,
		cancelFunc,
		connectionCancelFunc,
		logger,
		userContext,
		self.goFunctionCounter,
	)
}

func NewFactory(
	goFunctionCounter GoFunctionCounter.IService,
) *factory {
	return &factory{
		goFunctionCounter: goFunctionCounter,
	}
}
