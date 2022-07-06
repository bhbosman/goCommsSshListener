package channelListener

import (
	"github.com/bhbosman/goCommsSshListener/common"
	"sync"
)

type channelWithSemaphoreWrapper struct {
	common.IChannel
	mutex            sync.Mutex
	closed           bool
	releaseSemaphore interface {
		Release(int64)
	}
}

func (self *channelWithSemaphoreWrapper) Close() error {
	if self.closed {
		return nil
	}
	self.mutex.Lock()
	defer self.mutex.Unlock()
	if self.closed {
		return nil
	}
	self.closed = true
	err := self.IChannel.Close()
	local := self.releaseSemaphore
	self.releaseSemaphore = nil
	if local != nil {
		local.Release(1)
	}
	return err
}

func newSshChannelWithSemaphoreWrapper(
	conn common.IChannel,
	releaseSemaphore interface {
		Release(int64)
	}) *channelWithSemaphoreWrapper {
	return &channelWithSemaphoreWrapper{
		IChannel:         conn,
		mutex:            sync.Mutex{},
		releaseSemaphore: releaseSemaphore,
	}
}
