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

type SshChannel struct {
	channel                          ssh.Channel
	SpecificInformationForConnection goCommsDefinitions.ISpecificInformationForConnection
	name                             string
	additionalData                   []byte
}

func (self *SshChannel) ChannelType() string {
	return self.name
}

func (self *SshChannel) Read(data []byte) (int, error) {
	return self.channel.Read(data)
}

func (self *SshChannel) Write(data []byte) (int, error) {
	return self.channel.Write(data)
}

func (self *SshChannel) Close() error {
	return self.channel.Close()
}

func (self *SshChannel) CloseWrite() error {
	return self.channel.CloseWrite()
}

func (self *SshChannel) SendRequest(name string, wantReply bool, payload []byte) (bool, error) {
	return self.channel.SendRequest(name, wantReply, payload)
}

func (self *SshChannel) Stderr() io.ReadWriter {
	return self.channel.Stderr()
}

func (self *SshChannel) LocalAddr() net.Addr {
	return self.SpecificInformationForConnection.LocalAddr()
}

func (self *SshChannel) RemoteAddr() net.Addr {
	return self.SpecificInformationForConnection.RemoteAddr()
}

func (self *SshChannel) SetDeadline(t time.Time) error {
	return self.SpecificInformationForConnection.SetDeadline(t)
}

func (self *SshChannel) SetReadDeadline(t time.Time) error {
	return self.SpecificInformationForConnection.SetReadDeadline(t)
}

func (self *SshChannel) SetWriteDeadline(t time.Time) error {
	return self.SpecificInformationForConnection.SetWriteDeadline(t)
}

func NewSshChannel(
	name string,
	channel ssh.Channel,
	specificInformationForConnection goCommsDefinitions.ISpecificInformationForConnection,
	additionalData []byte,
) (common.ISshChannel, error) {
	if channel == nil {
		return nil, goerrors.InvalidParam
	}
	if specificInformationForConnection == nil {
		return nil, goerrors.InvalidParam
	}
	return &SshChannel{
		channel:                          channel,
		SpecificInformationForConnection: specificInformationForConnection,
		name:                             name,
		additionalData:                   additionalData,
	}, nil
}
