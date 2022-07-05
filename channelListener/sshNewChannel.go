package channelListener

import (
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goCommsSshListener/common"
	"github.com/bhbosman/goerrors"
	"golang.org/x/crypto/ssh"
	"net"
	"time"
)

type SshNewChannel struct {
	newChannel                       ssh.NewChannel
	specificInformationForConnection goCommsDefinitions.ISpecificInformationForConnection
}

func (self *SshNewChannel) LocalAddr() net.Addr {
	return self.specificInformationForConnection.LocalAddr()
}

func (self *SshNewChannel) RemoteAddr() net.Addr {
	return self.specificInformationForConnection.RemoteAddr()
}

func (self *SshNewChannel) SetDeadline(t time.Time) error {
	return self.specificInformationForConnection.SetDeadline(t)
}

func (self *SshNewChannel) SetReadDeadline(t time.Time) error {
	return self.specificInformationForConnection.SetReadDeadline(t)
}

func (self *SshNewChannel) SetWriteDeadline(t time.Time) error {
	return self.specificInformationForConnection.SetWriteDeadline(t)
}

func (self *SshNewChannel) Accept(
	channelType string,
	additionalData []byte,
) (common.ISshChannel, <-chan *ssh.Request, error) {
	accept, requests, err := self.newChannel.Accept()
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

func (self *SshNewChannel) Reject(reason ssh.RejectionReason, message string) error {
	return self.newChannel.Reject(reason, message)
}

func (self *SshNewChannel) ChannelType() string {
	return self.newChannel.ChannelType()
}

func (self *SshNewChannel) ExtraData() []byte {
	return self.newChannel.ExtraData()
}

func NewSshNewChannel(
	newChannel ssh.NewChannel,
	specificInformationForConnection goCommsDefinitions.ISpecificInformationForConnection,
) (common.ISshNewChannel, error) {
	if newChannel == nil {
		return nil, goerrors.InvalidParam
	}
	if specificInformationForConnection == nil {
		return nil, goerrors.InvalidParam
	}
	return &SshNewChannel{
		newChannel:                       newChannel,
		specificInformationForConnection: specificInformationForConnection,
	}, nil
}
