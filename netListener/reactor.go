package netListener

import (
	"context"
	internal2 "github.com/bhbosman/goCommsSshListener/internal"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/goprotoextra"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"sync/atomic"
)

type sshConnectionReactor struct {
	openCloserCount                int64
	onSend                         goprotoextra.ToConnectionFunc
	toConnectionReactor            goprotoextra.ToReactorFunc
	cancelCtx                      context.Context
	cancelFunc                     context.CancelFunc
	connectionCancelFunc           model.ConnectionCancelFunc
	logger                         *zap.Logger
	userContext                    interface{}
	messageHandlerService          IConnectionReactorMessageQueueService
	toConnectionFuncReplacement    rxgo.NextFunc
	toConnectionReactorReplacement rxgo.NextFunc
}

func (self *sshConnectionReactor) AddAcceptedChannel(uniqueReference string, acceptedChannel *fx.App) error {
	return self.messageHandlerService.AddAcceptedChannel(uniqueReference, acceptedChannel)
}

func (self *sshConnectionReactor) RemoveAcceptedChannel(uniqueReference string) error {
	return self.messageHandlerService.RemoveAcceptedChannel(uniqueReference)
}

func (self *sshConnectionReactor) CanAcceptChannel(name string) (bool, ssh.RejectionReason, string, error) {
	return self.messageHandlerService.CanAcceptChannel(name)
}

func (self *sshConnectionReactor) Close() error {
	addInt64 := atomic.AddInt64(&self.openCloserCount, -1)
	if addInt64 == 0 {
		return self.messageHandlerService.OnStop(context.Background())
	}
	return nil
}

func (self *sshConnectionReactor) Init(
	onSend goprotoextra.ToConnectionFunc,
	toConnectionReactor goprotoextra.ToReactorFunc,
	toConnectionFuncReplacement rxgo.NextFunc,
	toConnectionReactorReplacement rxgo.NextFunc,
) (rxgo.NextFunc, rxgo.ErrFunc, rxgo.CompletedFunc, chan interface{}, error) {
	self.onSend = onSend
	self.toConnectionReactor = toConnectionReactor
	self.toConnectionReactorReplacement = toConnectionReactorReplacement
	self.toConnectionFuncReplacement = toConnectionFuncReplacement
	return func(i interface{}) {
			_ = self.messageHandlerService.Send(i)
		},
		func(err error) {
			_ = self.messageHandlerService.Send(err)
		},
		func() {

		}, nil, nil
}

func (self *sshConnectionReactor) Open() error {
	addInt64 := atomic.AddInt64(&self.openCloserCount, 1)
	if addInt64 == 1 {
		return self.messageHandlerService.OnStart(context.Background())
	}
	return nil
}

func NewSshConnectionReactor(
	cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	connectionCancelFunc model.ConnectionCancelFunc,
	logger *zap.Logger,
	userContext interface{},
	goFunctionCounter GoFunctionCounter.IService,
) (internal2.ISshConnectionReactor, error) {

	onData := func() (IConnectionReactorMessageQueueData, error) {
		return NewData()
	}
	messageHandlerService, err := newService(
		cancelCtx,
		onData,
		logger,
		goFunctionCounter,
	)
	if err != nil {
		return nil, err
	}

	return &sshConnectionReactor{
		cancelCtx:             cancelCtx,
		cancelFunc:            cancelFunc,
		connectionCancelFunc:  connectionCancelFunc,
		logger:                logger,
		userContext:           userContext,
		messageHandlerService: messageHandlerService,
	}, nil
}
