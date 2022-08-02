package channelListener

import (
	"context"
	"github.com/bhbosman/goCommsSshListener/common"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/gdamore/tcell/v2/terminfo"
	"go.uber.org/multierr"
	"time"
)

type defaultChannelProcess struct {
	SetSizeCalled        bool
	Rows                 int
	Cols                 int
	sshChannel           common.IChannel
	timeOutCancelContext context.Context
	timeOutCancelFunc    context.CancelFunc
	closeCancelContext   context.Context
	closeCancelFunc      context.CancelFunc
	parentCancelFunc     context.CancelFunc
	Terminfo             *terminfo.Terminfo
	goFunctionCounter    GoFunctionCounter.IService
}

func (self *defaultChannelProcess) SetLookupTerminfo(terminfo *terminfo.Terminfo) error {
	self.Terminfo = terminfo
	return nil
}

func newDefaultChannelProcess(
	parentCancelContext context.Context,
	parentCancelFunc context.CancelFunc,
	sshChannel common.IChannel,
	goFunctionCounter GoFunctionCounter.IService,
) common.IChannelProcess {
	timeOutCancelContext, timeOutCancelFunc := context.WithTimeout(parentCancelContext, time.Second*100)
	closeCancelContext, closeCancelFunc := context.WithCancel(parentCancelContext)

	return &defaultChannelProcess{
		sshChannel:           sshChannel,
		timeOutCancelContext: timeOutCancelContext,
		timeOutCancelFunc:    timeOutCancelFunc,
		closeCancelContext:   closeCancelContext,
		closeCancelFunc:      closeCancelFunc,
		parentCancelFunc:     parentCancelFunc,
		goFunctionCounter:    goFunctionCounter,
	}
}

func (self *defaultChannelProcess) RunHandler() error {
	return self.goFunctionCounter.GoRun(
		"defaultChannelProcess.RunHandler",
		func() {
		loop:
			for {
				select {
				case <-self.timeOutCancelContext.Done():
					var errList error
					_, err := self.sshChannel.Stderr().Write([]byte("could not start remote process\n\r"))
					errList = multierr.Append(errList, err)

					errList = multierr.Append(errList, self.sshChannel.Close())
					if errList != nil {
						// logging
					}
					self.parentCancelFunc()
					break loop
				case <-self.closeCancelContext.Done():
					break loop
				}
			}
			self.timeOutCancelFunc()
		},
	)
}

func (self *defaultChannelProcess) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (self *defaultChannelProcess) Close() error {
	self.closeCancelFunc()

	return nil
}

func (self *defaultChannelProcess) ReadLine() (string, error) {
	return "", nil
}

func (self *defaultChannelProcess) SetSize(cols int, rows int) error {
	self.Cols = cols
	self.Rows = rows
	self.SetSizeCalled = true
	return nil
}
