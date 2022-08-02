package stack

import (
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/goprotoextra"
)

type OutboundStackHandler struct {
	errorState error
	stackData  *data
}

func (self *OutboundStackHandler) GetAdditionalBytesIncoming() int {
	return 0
}

func (self *OutboundStackHandler) GetAdditionalBytesSend() int {
	return self.stackData.connWrapper.BytesWritten
}

func (self *OutboundStackHandler) ReadMessage(_ interface{}) (interface{}, bool, error) {
	return nil, false, nil
}

func (self *OutboundStackHandler) Close() error {
	return self.stackData.Close()
}

func (self *OutboundStackHandler) SendRws(rws goprotoextra.ReadWriterSize) error {
	if self.errorState != nil {
		return self.errorState
	}
	return goerrors.NotImplemented
}

func (self *OutboundStackHandler) OnError(err error) {
	self.errorState = err
}

func (self *OutboundStackHandler) NextReadWriterSize(
	rws goprotoextra.ReadWriterSize,
	_ func(rws goprotoextra.ReadWriterSize) error,
	_ func(interface{}) error,
	_ func(size int) error) error {
	if self.errorState != nil {
		return self.errorState
	}
	return self.SendRws(rws)
}

func (self *OutboundStackHandler) OnComplete() {
	if self.errorState != nil {
		self.errorState = RxHandlers.RxHandlerComplete
	}
}

func NewOutboundStackHandler(stackData *data) (*OutboundStackHandler, error) {
	return &OutboundStackHandler{
		stackData: stackData,
	}, nil
}
