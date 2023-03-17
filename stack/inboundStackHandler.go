package stack

import (
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/goprotoextra"
	"io"
)

type InboundStackHandler struct {
	errorState error
	stackData  *data
}

func (self *InboundStackHandler) EmptyQueue() {
}

func (self *InboundStackHandler) ClearCounters() {
}

func (self *InboundStackHandler) PublishCounters(counters *model.PublishRxHandlerCounters) {
}

func (self *InboundStackHandler) GetAdditionalBytesIncoming() int {
	return 0
}

func (self *InboundStackHandler) GetAdditionalBytesSend() int {
	return 0
}

func (self *InboundStackHandler) ReadMessage(_ interface{}) error {
	return nil
}

func (self *InboundStackHandler) Close() error {
	return self.stackData.Close()
}

func (self *InboundStackHandler) SendRws(rws goprotoextra.ReadWriterSize) {
	_, err := io.Copy(self.stackData.pipeWriteClose, rws)
	if err != nil {
		self.errorState = err
		return
	}
}

func (self *InboundStackHandler) OnError(err error) {
	self.errorState = err
}

func (self *InboundStackHandler) NextReadWriterSize(
	rws goprotoextra.ReadWriterSize,
	_ func(rws goprotoextra.ReadWriterSize) error,
	_ func(interface{}) error,
	_ func(size int) error) error {

	if self.errorState != nil {
		return self.errorState
	}
	self.SendRws(rws)
	return self.errorState
}

func (self *InboundStackHandler) OnComplete() {
	if self.errorState != nil {
		self.errorState = RxHandlers.RxHandlerComplete
	}
}

func (self *InboundStackHandler) Complete() {

}

func NewInboundStackHandler(stackData *data) (*InboundStackHandler, error) {
	return &InboundStackHandler{
		stackData: stackData,
	}, nil
}
