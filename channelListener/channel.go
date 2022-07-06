package channelListener

import (
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goCommsSshListener/common"
	"github.com/bhbosman/goerrors"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"time"
)

type channel struct {
	sshChannel                       ssh.Channel
	SpecificInformationForConnection goCommsDefinitions.ISpecificInformationForConnection
	name                             string
	additionalData                   []byte
}

func (self *channel) ChannelType() string {
	return self.name
}

func (self *channel) Read(data []byte) (int, error) {
	return self.sshChannel.Read(data)
}

func (self *channel) Write(data []byte) (int, error) {
	return self.sshChannel.Write(data)
}

func (self *channel) Close() error {
	return self.sshChannel.Close()
}

func (self *channel) CloseWrite() error {
	return self.sshChannel.CloseWrite()
}

func (self *channel) SendRequest(name string, wantReply bool, payload []byte) (bool, error) {
	return self.sshChannel.SendRequest(name, wantReply, payload)
}

func (self *channel) Stderr() io.ReadWriter {
	return self.sshChannel.Stderr()
}

func (self *channel) LocalAddr() net.Addr {
	return self.SpecificInformationForConnection.LocalAddr()
}

func (self *channel) RemoteAddr() net.Addr {
	return self.SpecificInformationForConnection.RemoteAddr()
}

func (self *channel) SetDeadline(t time.Time) error {
	return self.SpecificInformationForConnection.SetDeadline(t)
}

func (self *channel) SetReadDeadline(t time.Time) error {
	return self.SpecificInformationForConnection.SetReadDeadline(t)
}

func (self *channel) SetWriteDeadline(t time.Time) error {
	return self.SpecificInformationForConnection.SetWriteDeadline(t)
}

func NewSshChannel(
	name string,
	sshChannel ssh.Channel,
	specificInformationForConnection goCommsDefinitions.ISpecificInformationForConnection,
	additionalData []byte,
) (common.IChannel, error) {
	if sshChannel == nil {
		return nil, goerrors.InvalidParam
	}
	if specificInformationForConnection == nil {
		return nil, goerrors.InvalidParam
	}
	return &channel{
		sshChannel:                       sshChannel,
		SpecificInformationForConnection: specificInformationForConnection,
		name:                             name,
		additionalData:                   additionalData,
	}, nil
}
