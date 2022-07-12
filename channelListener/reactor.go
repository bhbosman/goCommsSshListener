package channelListener

import (
	"context"
	"github.com/bhbosman/goCommsSshListener/common"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/messageRouter"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"github.com/gdamore/tcell/v2/terminfo"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"io"
)

type reactor struct {
	channelType                    string
	cancelCtx                      context.Context
	cancelFunc                     context.CancelFunc
	sshChannel                     common.IChannel
	toConnectionReactor            goprotoextra.ToReactorFunc
	channelProcess                 common.IChannelProcess
	messageRouter                  *messageRouter.MessageRouter
	onSend                         goprotoextra.ToConnectionFunc
	logger                         *zap.Logger
	sshChannelSessionSettings      common.ISshChannelSessionSettings
	extraData                      []byte
	goFunctionCounter              GoFunctionCounter.IService
	toConnectionFuncReplacement    rxgo.NextFunc
	toConnectionReactorReplacement rxgo.NextFunc
}

func (self *reactor) handleEmptyQueue(_ *messages.EmptyQueue) error {
	return nil
}

func (self *reactor) setProcess(newProcess common.IChannelProcess) error {
	old := self.channelProcess
	self.channelProcess = newProcess
	if defaultProcess, ok := old.(*defaultChannelProcess); ok && defaultProcess.SetSizeCalled {
		err := self.channelProcess.SetSize(defaultProcess.Cols, defaultProcess.Rows)
		if err != nil {
			return err
		}
		err = self.channelProcess.SetLookupTerminfo(defaultProcess.Terminfo)
		if err != nil {
			return err
		}
	}
	err := old.Close()
	if err != nil {
		return err
	}
	return nil
}

func (self *reactor) handleSshRequest(request *ssh.Request) {

	switch request.Type {
	//case "exec":
	//	msg := &execMsg{}
	//	err := ssh.Unmarshal(request.Payload, msg)
	//	if err != nil {
	//		_, _ = self.sshChannel.Stderr().Write([]byte(err.Error()))
	//		_ = self.sshChannel.Close()
	//	}
	//	cmd := exec.Command(msg.Command)
	//	cmd.StdinPipe() outPipe()
	//	cmd.Stdin = self.sshChannel
	//	cmd.Stderr = self.sshChannel.Stderr()
	//
	//	err = cmd.Run()
	//	if err != nil {
	//		_, _ = self.sshChannel.Stderr().Write([]byte(err.Error()))
	//		_ = self.sshChannel.Close()
	//	}

	case "shell":
		whatChannelToUse, err := self.sshChannelSessionSettings.UseDefault()
		if err != nil {
			err = self.sshChannel.Close()
			if err != nil {
				// TODO: Add informational logging.
			}
		}
		switch whatChannelToUse {
		case common.NoneInChannelProcess:
			break
		case common.EchoBuildInChannelProcess:
			err = self.createEchoChannelProcess()
			if err != nil {
				err = self.sshChannel.Close()
				if err != nil {
					// TODO: Add informational logging.
				}
			}
			break
		case common.UseCustomChannelProcess:
			var newProcess common.IChannelProcess
			newProcess, err = self.sshChannelSessionSettings.CreateChannelProcess(
				self.sshChannel,
				self.cancelCtx,
				self.cancelFunc,
				self.onSend,
				self.toConnectionFuncReplacement,
				self.logger,
			)
			if err != nil {
				err = self.sshChannel.Close()
				if err != nil {
					// TODO: Add informational logging.
				}
			}
			err = self.assignChannelProcess(newProcess)
			if err != nil {
				err = self.sshChannel.Close()
				if err != nil {
					// TODO: Add informational logging.
				}
			}
			break
		}
		break
	case "pty-req":

		msg := &ptyRequestMsg{}
		err := ssh.Unmarshal(request.Payload, msg)
		if err != nil {
		}
		err = self.channelProcess.SetSize(int(msg.Columns), int(msg.Rows))
		if err != nil {
		}
		lookupTerminfo, err := terminfo.LookupTerminfo(msg.Term)
		if err != nil {
		}
		err = self.channelProcess.SetLookupTerminfo(lookupTerminfo)
		if err != nil {
		}
		break
	case "window-change":
		msg := &ptyWindowChangeMsg{}
		_ = ssh.Unmarshal(request.Payload, msg)
		_ = self.channelProcess.SetSize(int(msg.Columns), int(msg.Rows))
		break
	}
	if request.WantReply {
		_ = request.Reply(true, nil)
	}
}

func (self *reactor) handleReaderWriter(message *gomessageblock.ReaderWriter) {
	_, _ = io.Copy(self.channelProcess, message)
}

func (self *reactor) Init(
	onSend goprotoextra.ToConnectionFunc,
	toConnectionReactor goprotoextra.ToReactorFunc,
	toConnectionFuncReplacement rxgo.NextFunc,
	toConnectionReactorReplacement rxgo.NextFunc,
) (rxgo.NextFunc, rxgo.ErrFunc, rxgo.CompletedFunc, chan interface{}, error) {
	self.toConnectionReactor = toConnectionReactor
	self.onSend = onSend
	self.toConnectionFuncReplacement = toConnectionFuncReplacement
	self.toConnectionReactorReplacement = toConnectionReactorReplacement

	return func(i interface{}) {
			_, _ = self.messageRouter.Route(i)
		},
		func(err error) {
			_, _ = self.messageRouter.Route(err)
		},
		func() {

		}, nil, nil
}

func (self *reactor) Open() error {
	return nil
}

func (self *reactor) Close() error {
	var err error
	if self.channelProcess != nil {
		err = multierr.Append(err, self.channelProcess.Close())
	}
	return err
}

func (self *reactor) createEchoChannelProcess() error {
	newProcess, err := newEchoShellProcess(
		self.sshChannel,
		self.cancelCtx,
		self.cancelFunc,
		self.onSend,
		self.toConnectionFuncReplacement,
		self.goFunctionCounter,
	)
	if err != nil {
		self.logger.Error("On CreateEchoChannelProcess/newEchoShellProcess", zap.Error(err))
		err = multierr.Append(err, self.sshChannel.Close())
		if err != nil {
			self.logger.Error("On CreateEchoChannelProcess channel close", zap.Error(err))
		}
		return err
	}
	return self.assignChannelProcess(newProcess)
}

func (self *reactor) assignChannelProcess(process common.IChannelProcess) error {

	err := self.setProcess(process)
	if err != nil {
		self.logger.Error("On CreateEchoChannelProcess/setProcess", zap.Error(err))
		err = multierr.Append(err, self.sshChannel.Close())
		if err != nil {
			self.logger.Error("On CreateEchoChannelProcess channel close", zap.Error(err))
		}
		return err
	}

	err = process.RunHandler()
	if err != nil {
		self.logger.Error("On CreateEchoChannelProcess/RunHandler", zap.Error(err))
		err = multierr.Append(err, self.sshChannel.Close())
		if err != nil {
			self.logger.Error("On CreateEchoChannelProcess channel close", zap.Error(err))
		}
		return err
	}
	return nil
}

func NewReactor(
	channelType string,
	cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	logger *zap.Logger,
	sshChannel common.IChannel,
	sshChannelSessionSettings common.ISshChannelSessionSettings,
	extraData []byte,
	goFunctionCounter GoFunctionCounter.IService,
) (intf.IConnectionReactor, error) {
	defaultSessionProcess := newDefaultChannelProcess(
		cancelCtx,
		cancelFunc,
		sshChannel,
		goFunctionCounter,
	)
	err := defaultSessionProcess.RunHandler()
	if err != nil {
		return nil, err
	}

	reactorInstance := &reactor{
		channelType:               channelType,
		cancelCtx:                 cancelCtx,
		cancelFunc:                cancelFunc,
		sshChannel:                sshChannel,
		channelProcess:            defaultSessionProcess,
		messageRouter:             messageRouter.NewMessageRouter(),
		logger:                    logger,
		sshChannelSessionSettings: sshChannelSessionSettings,
		extraData:                 extraData,
		goFunctionCounter:         goFunctionCounter,
	}

	reactorInstance.messageRouter.Add(reactorInstance.handleEmptyQueue)
	reactorInstance.messageRouter.Add(reactorInstance.handleSshRequest)
	reactorInstance.messageRouter.Add(reactorInstance.handleReaderWriter)

	return reactorInstance, nil
}

type ptyRequestMsg struct {
	Term     string
	Columns  uint32
	Rows     uint32
	Width    uint32
	Height   uint32
	Modelist string
}

type execMsg struct {
	Command string
}

type ptyWindowChangeMsg struct {
	Columns uint32
	Rows    uint32
	Width   uint32
	Height  uint32
}
