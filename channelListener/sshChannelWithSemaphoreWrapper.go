package channelListener

import (
	"github.com/bhbosman/goCommsSshListener/common"
	"sync"
)

type sshChannelWithSemaphoreWrapper struct {
	common.ISshChannel
	mutex            sync.Mutex
	closed           bool
	releaseSemaphore interface {
		Release(int64)
	}
}

func (self *sshChannelWithSemaphoreWrapper) Close() error {
	if self.closed {
		return nil
	}
	self.mutex.Lock()
	defer self.mutex.Unlock()
	if self.closed {
		return nil
	}
	self.closed = true
	err := self.ISshChannel.Close()
	local := self.releaseSemaphore
	self.releaseSemaphore = nil
	if local != nil {
		local.Release(1)
	}
	return err
}

func newSshChannelWithSemaphoreWrapper(
	conn common.ISshChannel,
	releaseSemaphore interface {
		Release(int64)
	}) *sshChannelWithSemaphoreWrapper {
	return &sshChannelWithSemaphoreWrapper{
		ISshChannel:      conn,
		mutex:            sync.Mutex{},
		releaseSemaphore: releaseSemaphore,
	}
}
