package channelListener

import (
	"context"
	"github.com/bhbosman/goCommsSshListener/common"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/reactivex/rxgo/v2"
	"golang.org/x/crypto/ssh/terminal"
)

type echoShellProcess struct {
	common.BaseChannelProcess
	handler           *terminal.Terminal
	goFunctionCounter GoFunctionCounter.IService
}

func newEchoShellProcess(
	sshChannel common.IChannel,
	parentCtx context.Context,
	parentCancelFunc context.CancelFunc,
	onSend rxgo.NextFunc,
	goFunctionCounter GoFunctionCounter.IService,
) (*echoShellProcess, error) {
	emptyShell := common.NewBaseChannelProcess(
		sshChannel,
		parentCtx,
		parentCancelFunc,
		onSend,
	)
	newProcess := terminal.NewTerminal(emptyShell.RwProxy, ">>")
	return &echoShellProcess{
		BaseChannelProcess: emptyShell,
		handler:            newProcess,
		goFunctionCounter:  goFunctionCounter,
	}, nil
}

func (self *echoShellProcess) RunHandler() error {
	return self.goFunctionCounter.GoRun(
		"echoShellProcess.RunHandler",
		func() {
			for self.CancelCtx.Err() == nil {
				line, err := self.handler.ReadLine()
				if err != nil {
					self.CancelFunc()
				}
				if line != "" {
				}
			}
			_ = self.SshChannel.Close()
		},
	)

}

func (self *echoShellProcess) SetSize(cols int, rows int) error {
	return self.handler.SetSize(cols, rows)
}

func (self *echoShellProcess) ReadLine() (string, error) {
	return self.handler.ReadLine()
}
