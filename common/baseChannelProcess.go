package common

import (
	"context"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/goprotoextra"
	"github.com/gdamore/tcell/v2/terminfo"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/multierr"
	"io"
)

type BaseChannelProcess struct {
	CancelCtx      context.Context
	CancelFunc     context.CancelFunc
	SshChannel     ISshChannel
	RwProxy        io.ReadWriteCloser
	PipeWriteClose io.WriteCloser
	Terminfo       *terminfo.Terminfo
}

func NewBaseChannelProcess(
	sshChannel ISshChannel,
	parentCtx context.Context,
	parentCancelFunc context.CancelFunc,
	onSend goprotoextra.ToConnectionFunc,
	onSendReplacement rxgo.NextFunc,
) BaseChannelProcess {
	tempPipeReadClose, tempPipeWriteClose := common.Pipe(parentCtx)
	rwProxy := &ReaderWriterProxy{
		PipeReader:        tempPipeReadClose,
		OnSend:            onSend,
		onSendReplacement: onSendReplacement,
	}
	return BaseChannelProcess{
		CancelCtx:      parentCtx,
		CancelFunc:     parentCancelFunc,
		SshChannel:     sshChannel,
		RwProxy:        rwProxy,
		PipeWriteClose: tempPipeWriteClose,
	}
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
