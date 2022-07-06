package channelListener

import (
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goCommsSshListener/common"
	"golang.org/x/crypto/ssh"
	"golang.org/x/net/context"
	"io"
	"net"
	url2 "net/url"
)

type iChannelListener interface {
	io.Closer
	accept() (common.INewChannel, error)
	addr() net.Addr
}

type channelListener struct {
	cancelContext  context.Context
	cancelFunc     context.CancelFunc
	channel        <-chan ssh.NewChannel
	Url            *url2.URL
	address        net.Addr
	connectionInfo goCommsDefinitions.ISpecificInformationForConnection
}

func (self *channelListener) accept() (common.INewChannel, error) {
	select {
	case ch, ok := <-self.channel:
		if ok {
			return NewSshNewChannel(
				ch,
				self.connectionInfo)
		}
		return nil, io.EOF
	case <-self.cancelContext.Done():
		return nil, io.EOF
	}
}

func (self *channelListener) Close() error {
	return nil
}

func (self *channelListener) addr() net.Addr {
	return self.address
}

func newChannelListener(
	parentContext context.Context,
	channel <-chan ssh.NewChannel,
	SpecificInformationForConnection goCommsDefinitions.ISpecificInformationForConnection,
	Url *url2.URL,
) (iChannelListener, error) {
	cancelContext, cancelFunc := context.WithCancel(parentContext)
	return &channelListener{
		cancelContext: cancelContext,
		cancelFunc:    cancelFunc,
		channel:       channel,
		Url:           Url,
		address: &sshChannelAddress{
			path: Url,
		},
		connectionInfo: SpecificInformationForConnection,
	}, nil
}

type sshChannelAddress struct {
	path *url2.URL
}

func (self *sshChannelAddress) Network() string {
	return self.path.Scheme
}

func (self *sshChannelAddress) String() string {
	return self.path.Path
}
