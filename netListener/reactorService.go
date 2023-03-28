package netListener

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/gocommon/ChannelHandler"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocommon/services/IFxService"
	"github.com/bhbosman/gocommon/services/ISendMessage"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
)

type service struct {
	ctx               context.Context
	cancelFunc        context.CancelFunc
	onData            func() (IConnectionReactorMessageQueueData, error)
	logger            *zap.Logger
	state             IFxService.State
	cmdChannel        chan interface{}
	goFunctionCounter GoFunctionCounter.IService
}

func (self *service) Send(message interface{}) error {
	result, err := ISendMessage.CallISendMessageSend(
		self.ctx, self.cmdChannel, true,
		message)
	if err != nil {
		return err
	}
	return result.Args0
}

func (self *service) AddAcceptedChannel(uniqueReference string, acceptedChannel messages.IApp) error {
	if self.state == IFxService.Started {
		result, err := CallIConnectionReactorMessageQueueAddAcceptedChannel(
			self.ctx, self.cmdChannel, true,
			uniqueReference, acceptedChannel)
		if err != nil {
			return err
		}
		return result.Args0
	}
	return nil
}

func (self *service) RemoveAcceptedChannel(uniqueReference string) error {
	result, err := CallIConnectionReactorMessageQueueRemoveAcceptedChannel(
		self.ctx, self.cmdChannel, true,
		uniqueReference)
	if err != nil {
		return err
	}
	return result.Args0
}

func (self *service) HandleGlobalRequest(name string, req *ssh.Request) error {
	request, err := CallIConnectionReactorMessageQueueHandleGlobalRequest(
		self.ctx, self.cmdChannel, true, name, req)
	if err != nil {
		return err
	}
	return request.Args0
}

func (self *service) CanAcceptChannel(name string) (bool, ssh.RejectionReason, string, error) {
	request, err := CallIConnectionReactorMessageQueueCanAcceptChannel(
		self.ctx, self.cmdChannel, true, name)
	if err != nil {
		return false, ssh.ResourceShortage, "", err
	}
	return request.Args0, request.Args1, request.Args2, request.Args3
}

func (self *service) OnStart(_ context.Context) error {
	err := self.start()
	if err != nil {
		return err
	}
	self.state = IFxService.Started
	return nil
}

func (self *service) OnStop(_ context.Context) error {
	err := self.shutdown()
	close(self.cmdChannel)
	self.state = IFxService.Stopped
	return err
}

func (self *service) State() IFxService.State {
	return self.state
}

func (self *service) ServiceName() string {
	return "ConnectionReactorMessageHandler"
}

func (self *service) start() error {
	data, err := self.onData()
	if err != nil {
		return err
	}

	return self.goFunctionCounter.GoRun(
		"IConnectionReactorMessageQueueService.Start",
		func() {
			self.goStart(data)
		},
	)
}

func (self *service) goStart(data IConnectionReactorMessageQueueData) {
	defer func(cmdChannel <-chan interface{}) {
		//flush
		for range cmdChannel {
		}
	}(self.cmdChannel)
	channelHandlerCallback := ChannelHandler.CreateChannelHandlerCallback(
		self.ctx,
		data,
		[]ChannelHandler.ChannelHandler{
			{
				Cb: func(next interface{}, message interface{}) (bool, error) {
					if unk, ok := next.(IConnectionReactorMessageQueue); ok {
						return ChannelEventsForIConnectionReactorMessageQueue(unk, message)
					}
					return false, nil
				},
			},
			{
				Cb: func(next interface{}, message interface{}) (bool, error) {
					if unk, ok := next.(ISendMessage.ISendMessage); ok {
						return ISendMessage.ChannelEventsForISendMessage(unk, message)
					}
					return false, nil
				},
			},
		},
		func() int {
			return len(self.cmdChannel)
		},
		goCommsDefinitions.CreateTryNextFunc(self.cmdChannel),
	)
loop:
	for {
		select {
		case <-self.ctx.Done():
			err := data.ShutDown()
			if err != nil {
				self.logger.Error(
					"error on done",
					zap.Error(err))
			}
			break loop
		case event, ok := <-self.cmdChannel:
			if !ok {
				return
			}
			breakLoop, err := channelHandlerCallback(event)
			if err != nil || breakLoop {
				break loop
			}
		}
	}
}

func (self *service) shutdown() error {
	self.cancelFunc()
	return nil
}

func newService(
	ctx context.Context,
	onData func() (IConnectionReactorMessageQueueData, error),
	logger *zap.Logger,
	goFunctionCounter GoFunctionCounter.IService,
) (IConnectionReactorMessageQueueService, error) {
	localCtx, localCancelFunc := context.WithCancel(ctx)

	return &service{
		ctx:               localCtx,
		cancelFunc:        localCancelFunc,
		onData:            onData,
		logger:            logger,
		cmdChannel:        make(chan interface{}, 32),
		goFunctionCounter: goFunctionCounter,
	}, nil
}
