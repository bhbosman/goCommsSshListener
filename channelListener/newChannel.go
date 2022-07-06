package channelListener

import (
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goCommsSshListener/common"
	"github.com/bhbosman/goerrors"
	"golang.org/x/crypto/ssh"
	"net"
	"time"
)

type newChannel struct {
	sshNewChannel                    ssh.NewChannel
	specificInformationForConnection goCommsDefinitions.ISpecificInformationForConnection
}

func (self *newChannel) LocalAddr() net.Addr {
	return self.specificInformationForConnection.LocalAddr()
}

func (self *newChannel) RemoteAddr() net.Addr {
	return self.specificInformationForConnection.RemoteAddr()
}

func (self *newChannel) SetDeadline(t time.Time) error {
	return self.specificInformationForConnection.SetDeadline(t)
}

func (self *newChannel) SetReadDeadline(t time.Time) error {
	return self.specificInformationForConnection.SetReadDeadline(t)
}

func (self *newChannel) SetWriteDeadline(t time.Time) error {
	return self.specificInformationForConnection.SetWriteDeadline(t)
}

func (self *newChannel) Accept(
	channelType string,
	additionalData []byte,
) (common.IChannel, <-chan *ssh.Request, error) {
	accept, requests, err := self.sshNewChannel.Accept()
	sshChannel, err := NewSshChannel(
		channelType,
		accept,
		self.specificInformationForConnection,
		additionalData,
	)
	if err != nil {
		return nil, nil, err
	}
	return sshChannel, requests, err
}

func (self *newChannel) Reject(reason ssh.RejectionReason, message string) error {
	return self.sshNewChannel.Reject(reason, message)
}

func (self *newChannel) ChannelType() string {
	return self.sshNewChannel.ChannelType()
}

func (self *newChannel) ExtraData() []byte {
	return self.sshNewChannel.ExtraData()
}

func NewSshNewChannel(
	sshNewChannel ssh.NewChannel,
	specificInformationForConnection goCommsDefinitions.ISpecificInformationForConnection,
) (common.INewChannel, error) {
	if sshNewChannel == nil {
		return nil, goerrors.InvalidParam
	}
	if specificInformationForConnection == nil {
		return nil, goerrors.InvalidParam
	}
	return &newChannel{
		sshNewChannel:                    sshNewChannel,
		specificInformationForConnection: specificInformationForConnection,
	}, nil
}
