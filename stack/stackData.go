package stack

import (
	"context"
	"fmt"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goCommsSshListener/channelListener"
	common2 "github.com/bhbosman/goCommsSshListener/common"
	"github.com/bhbosman/goCommsSshListener/internal"
	"github.com/bhbosman/goConnectionManager"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/services/interfaces"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/goerrors"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"net/url"
)

type data struct {
	connectionType      model.ConnectionType
	sshConn             ssh.Conn
	conn                goCommsDefinitions.ISpecificInformationForConnection
	ctx                 context.Context
	cancelFunc          context.CancelFunc
	logger              *zap.Logger
	pipeWriteClose      io.WriteCloser
	pipeRead            io.ReadCloser
	connWrapper         *common.ConnWrapper
	connectionManager   goConnectionManager.IService
	uniqueSessionNumber interfaces.IUniqueReferenceService
	connectionReactor   internal.ISshConnectionReactor
	sshChannelSettings  common2.ISshChannelSettings
	goFunctionCounter   GoFunctionCounter.IService
	outBoundHandler     goCommsDefinitions.IRxNextHandler
	inBoundHandler      goCommsDefinitions.IRxNextHandler
}

func (self *data) Close() error {
	var errList error = nil

	if self.sshConn != nil {
		errList = multierr.Append(errList, self.sshConn.Close())
	} else {
		errList = multierr.Append(errList, goerrors.InvalidParam)
	}

	if self.pipeWriteClose != nil {
		errList = multierr.Append(errList, self.pipeWriteClose.Close())
	} else {
		errList = multierr.Append(errList, goerrors.InvalidParam)
	}

	if self.connWrapper != nil {
		errList = multierr.Append(errList, self.connWrapper.Close())
	} else {
		errList = multierr.Append(errList, goerrors.InvalidParam)
	}

	return errList
}

func (self *data) Start(
	Ctx context.Context,
	ToReactorFunc rxgo.NextFunc,
) (common.IInputStreamForStack, error) {
	if self.connectionType == model.ClientConnection {
		return nil, goerrors.InvalidParam
	}
	config := &ssh.ServerConfig{
		Config: ssh.Config{
			Rand:           nil,
			RekeyThreshold: 0,
			KeyExchanges:   nil,
			Ciphers:        nil,
			MACs:           nil,
		},
		NoClientAuth:     false,
		MaxAuthTries:     0,
		PasswordCallback: nil,
		PublicKeyCallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			return &ssh.Permissions{
				Extensions: map[string]string{
					"key-id": "sssssss",
				},
			}, nil
		},
		KeyboardInteractiveCallback: nil,
		AuthLogCallback:             nil,
		ServerVersion:               "",
		BannerCallback:              nil,
		GSSAPIWithMICConfig:         nil,
	}

	key, err := ssh.ParsePrivateKey(_defaultcert)
	if err != nil {
		return nil, err
	}
	config.AddHostKey(key)

	serverConn, channels, requests, err := ssh.NewServerConn(self.connWrapper, config)
	if err != nil {
		return nil, err
	}

	err = self.handleChannelsAndRequest(self.conn, channels, requests, ToReactorFunc)
	if err != nil {
		return nil, err
	}
	self.sshConn = serverConn

	if err != nil {
		return nil, err
	}

	return nil, Ctx.Err()
}

func (self *data) Stop() error {
	return nil
}

func (self *data) handleChannelsAndRequest(
	fromConnection goCommsDefinitions.ISpecificInformationForConnection,
	channels <-chan ssh.NewChannel,
	globalRequests <-chan *ssh.Request,
	ToReactorFunc rxgo.NextFunc,
) error {

	netAppParams := common.NewNetAppFuncInParams(
		self.ctx,
		self.connectionManager,
		self.logger,
		self.uniqueSessionNumber,
		self.goFunctionCounter,
	)
	err := internal.GoRequestChannelHandler(
		"handleChannelsAndRequest",
		self.ctx,
		globalRequests,
		ToReactorFunc,
		self.goFunctionCounter,
	)
	if err != nil {
		return err
	}
	var parse *url.URL
	parse, err = url.Parse(
		fmt.Sprintf("sshChannelAcceptor:///?local_address=%v&remote_address=%v",
			fromConnection.LocalAddr().String(),
			fromConnection.RemoteAddr().String()))
	if err != nil {
		return err
	}

	app := channelListener.NewManagerApp(
		"sasdasdas",
		"asdsa",
		parse.String(),
		channels,
		self.conn,
		common.MoreOptions(
			fx.Provide(
				fx.Annotated{
					Target: func() internal.ISshConnectionReactor {
						return self.connectionReactor
					},
				},
			),
		),
		common.MoreOptions(
			fx.Provide(
				fx.Annotated{
					Target: func() common2.ISshChannelSettings {
						return self.sshChannelSettings
					},
				},
			),
		),
	)

	callback, cancelFunc, err := app(netAppParams).Callback()
	onError := func() {
		if cancelFunc != nil {
			//cancelFunc()
		}
	}

	if err != nil {
		onError()
		return err
	}
	if callback.Err() != nil {
		onError()
		return err
	}
	err = callback.Start(context.Background())
	if err != nil {
		onError()
		return err
	}

	return self.goFunctionCounter.GoRun(
		"SSHStackData.handleChannelsAndRequest",
		func() {
			<-self.ctx.Done()
			_ = callback.Stop(context.Background())
		},
	)
}

//
//func (self *data) setInBoundHandler(handlers *goCommsDefinitions.DefaultRxNextHandler) error {
//	if handlers == nil {
//		return goerrors.InvalidParam
//	}
//	self.inBoundHandler = handlers
//	return nil
//
//}

func (self *data) setConnWrapper(wrapper *common.ConnWrapper) error {
	self.connWrapper = wrapper
	return nil
}

func NewStackData(
	connectionType model.ConnectionType,
	conn net.Conn,
	ctx context.Context,
	cancelFunc context.CancelFunc,
	logger *zap.Logger,
	ConnectionManager goConnectionManager.IService,
	uniqueSessionNumber interfaces.IUniqueReferenceService,
	connectionReactor internal.ISshConnectionReactor,
	SshChannelSettings common2.ISshChannelSettings,
	goFunctionCounter GoFunctionCounter.IService,
) *data {
	tempPipeRead, tempPipeWriteClose := common.Pipe(ctx)

	return &data{
		connectionType:      connectionType,
		conn:                conn,
		ctx:                 ctx,
		cancelFunc:          cancelFunc,
		logger:              logger,
		pipeWriteClose:      tempPipeWriteClose,
		pipeRead:            tempPipeRead,
		connectionManager:   ConnectionManager,
		uniqueSessionNumber: uniqueSessionNumber,
		connectionReactor:   connectionReactor,
		sshChannelSettings:  SshChannelSettings,
		goFunctionCounter:   goFunctionCounter,
	}
}
