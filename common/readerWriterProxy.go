package common

import (
	"github.com/bhbosman/gomessageblock"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/multierr"
	"io"
)

type ReaderWriterProxy struct {
	OnSend     rxgo.NextFunc
	PipeReader io.ReadCloser
}

func (self *ReaderWriterProxy) Close() error {
	var err error
	if self.PipeReader != nil {
		err = multierr.Append(err, self.PipeReader.Close())
	}
	return err
}

func (self *ReaderWriterProxy) Read(p []byte) (n int, err error) {
	return self.PipeReader.Read(p)
}

func (self *ReaderWriterProxy) Write(p []byte) (n int, err error) {
	self.OnSend(gomessageblock.NewReaderWriterBlock(p))
	if err != nil {
		return 0, err
	}
	n = len(p)
	return n, err
}
