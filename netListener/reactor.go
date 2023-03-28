package netListener

import (
	"context"
	internal2 "github.com/bhbosman/goCommsSshListener/internal"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/services/IFxService"
	"github.com/bhbosman/gocomms/intf"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"sync/atomic"
)

type sshConnectionReactor struct {
	openCloserCount                int64
	onSendToReactor                rxgo.NextFunc
	onSendToConnection             rxgo.NextFunc
	onSend                         rxgo.NextFunc
	cancelCtx                      context.Context
	cancelFunc                     context.CancelFunc
	connectionCancelFunc           model.ConnectionCancelFunc
	logger                         *zap.Logger
	messageHandlerService          IConnectionReactorMessageQueueService
	toConnectionFuncReplacement    rxgo.NextFunc
	toConnectionReactorReplacement rxgo.NextFunc
}

func (self *sshConnectionReactor) AddAcceptedChannel(uniqueReference string, acceptedChannel messages.IApp) error {
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

func (self *sshConnectionReactor) Init(params intf.IInitParams) (rxgo.NextFunc, rxgo.ErrFunc, rxgo.CompletedFunc, error) {
	self.onSendToReactor = params.OnSendToReactor()
	self.onSendToConnection = params.OnSendToConnection()
	return func(i interface{}) {
			if self.messageHandlerService.State() == IFxService.Started {
				_ = self.messageHandlerService.Send(i)
			}
		},
		func(err error) {
			if self.messageHandlerService.State() == IFxService.Started {
				_ = self.messageHandlerService.Send(err)
			}
		},
		func() {

		}, nil
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
		cancelCtx:            cancelCtx,
		cancelFunc:           cancelFunc,
		connectionCancelFunc: connectionCancelFunc,
		logger:               logger,
		//userContext:           userContext,
		messageHandlerService: messageHandlerService,
	}, nil
}
