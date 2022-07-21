package common

import (
	"context"
	"github.com/bhbosman/gocomms/common"
	"github.com/gdamore/tcell/v2/terminfo"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/multierr"
	"io"
)

type BaseChannelProcess struct {
	CancelCtx      context.Context
	CancelFunc     context.CancelFunc
	SshChannel     IChannel
	RwProxy        io.ReadWriteCloser
	PipeWriteClose io.WriteCloser
	Terminfo       *terminfo.Terminfo
}

func (self *BaseChannelProcess) SetLookupTerminfo(terminfo *terminfo.Terminfo) error {
	self.Terminfo = terminfo
	return nil
}

func (self *BaseChannelProcess) Write(p []byte) (n int, err error) {
	return self.PipeWriteClose.Write(p)
}

func (self *BaseChannelProcess) Close() error {
	var err error
	if self.PipeWriteClose != nil {
		err = multierr.Append(err, self.PipeWriteClose.Close())
	}
	if self.RwProxy != nil {
		err = multierr.Append(err, self.RwProxy.Close())
	}
	return err
}

func NewBaseChannelProcess(
	sshChannel IChannel,
	parentCtx context.Context,
	parentCancelFunc context.CancelFunc,
	onSend rxgo.NextFunc,
) BaseChannelProcess {
	tempPipeReadClose, tempPipeWriteClose := common.Pipe(parentCtx)
	rwProxy := &ReaderWriterProxy{
		PipeReader: tempPipeReadClose,
		OnSend:     onSend,
	}
	return BaseChannelProcess{
		CancelCtx:      parentCtx,
		CancelFunc:     parentCancelFunc,
		SshChannel:     sshChannel,
		RwProxy:        rwProxy,
		PipeWriteClose: tempPipeWriteClose,
	}
}
