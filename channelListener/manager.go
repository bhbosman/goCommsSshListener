package channelListener

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goCommsSshListener/common"
	"github.com/bhbosman/goCommsSshListener/internal"
	"github.com/bhbosman/goConnectionManager"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/Services/IFxService"
	"github.com/bhbosman/gocommon/Services/interfaces"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/netBase"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"golang.org/x/sync/semaphore"
	"net/url"
)

type manager struct {
	netBase.ConnNetManager
	listener           ISshListenerAccept
	MaxConnections     int
	OnCreateConnection goCommsDefinitions.IOnCreateConnection
	ConnectionReactor  internal.ISshConnectionReactor
	SshChannelSettings common.ISshChannelSettings
}

func (self *manager) ListenForNewConnections() error {
	goFunc := func() {
		n := 0
		sem := semaphore.NewWeighted(int64(self.MaxConnections))
	loop:
		for self.CancelCtx.Err() == nil {
			n++
			self.ZapLogger.Info(
				"Trying to accept connections",
				zap.Int("Connection Count", n),
			)
			acceptNewChannel, connCancelFunc, err := self.acceptWithContext()
			if err != nil || err == nil && acceptNewChannel == nil {
				self.ZapLogger.Error(
					"Error on accept",
					zap.Error(err))
				continue loop
			}

			if sem.TryAcquire(1) {
				self.ZapLogger.Info("Accepted connection...")
				channelType := acceptNewChannel.ChannelType()
				var canAcceptChannel bool
				var rejectionReason ssh.RejectionReason
				var rejectionMessage string
				canAcceptChannel, rejectionReason, rejectionMessage, err = self.ConnectionReactor.CanAcceptChannel(channelType)

				if err != nil {
					self.ZapLogger.Error("Error on CanAcceptChannel", zap.Error(err))
					err = acceptNewChannel.Reject(ssh.Prohibited, "Error on asking CanAcceptChannel")
					if err != nil {
						self.ZapLogger.Error("On acceptNewChannel.Reject", zap.Error(err))
					}
					continue loop
				}

				if !canAcceptChannel {
					err = acceptNewChannel.Reject(rejectionReason, rejectionMessage)
					if err != nil {
						self.ZapLogger.Error("On acceptNewChannel.Reject when !canAcceptChannel", zap.Error(err))
					}
					continue loop
				}
				extraData := acceptNewChannel.ExtraData()
				var acceptedChannel common.ISshChannel
				var acceptedChannelRequestChannel <-chan *ssh.Request
				acceptedChannel, acceptedChannelRequestChannel, err = acceptNewChannel.Accept(
					channelType,
					extraData,
				)
				if err != nil {
					self.ZapLogger.Error("Error on acceptNewChannel.Accept", zap.Error(err))
					err = acceptNewChannel.Reject(ssh.Prohibited, "error on accepting channel")
					if err != nil {
						self.ZapLogger.Error("error on accepting channel", zap.Error(err))
					}
					continue loop
				}

				onErrorFlush := func(
					acceptedChannel common.ISshChannel,
					acceptedChannelRequestChannel <-chan *ssh.Request,
					cancelFunc context.CancelFunc, logger *zap.Logger) func(err error) {
					called := false
					return func(err error) {
						if !called {
							called = true
							cancelFunc()
							err = multierr.Append(err, acceptedChannel.Close())

							// this function is part of the GoFunctionCounter count
							go func() {
								functionName := self.GoFunctionCounter.CreateFunctionName("SshChannelListenerManager.ListenForNewConnections.Flush.AcceptedChannelRequestChannel")
								defer func(GoFunctionCounter GoFunctionCounter.IService, name string) {
									_ = GoFunctionCounter.Remove(name)
								}(self.GoFunctionCounter, functionName)
								_ = self.GoFunctionCounter.Add(functionName)

								//
								for range acceptedChannelRequestChannel {
									// do nothing
								}
							}()
							if err != nil {
								logger.Error("OnErrorFlush", zap.Error(err))
							}
						}
					}
				}(acceptedChannel, acceptedChannelRequestChannel, connCancelFunc, self.ZapLogger)

				acceptedChannel = newSshChannelWithSemaphoreWrapper(acceptedChannel, sem)

				// override the stack here to use the session's connection stack
				stackName := self.StackName
				switch channelType {
				case "session":
					stackName = goCommsDefinitions.TransportFactoryForSshChannelSession
				default:

				}
				uniqueReference := self.UniqueSessionNumber.Next(self.ConnectionInstancePrefix)
				connectionApp, ctx, cancelFunc := self.NewConnectionInstanceWithStackName(
					uniqueReference,
					self.GoFunctionCounter,
					model.ServerConnection,
					acceptedChannel,
					stackName,
					netBase.NewAddFxOptions(
						provideCreateIConnectionReactor(),
						fx.Provide(
							fx.Annotated{
								Name: "ChannelExtraData",
								Target: func() ([]byte, error) {
									return extraData, nil
								},
							},
						),
						fx.Provide(
							fx.Annotated{
								Target: func(
									params struct {
										fx.In
									},
								) (common.ISshChannel, error) {
									return acceptedChannel, nil
								},
							},
						),
						provideChannelType(channelType),
						provideSshConnectionReactor(self.ConnectionReactor),
						provideSshCreateChannelProcess(self.SshChannelSettings),
						provideAcceptedChannelRequestChannel(acceptedChannelRequestChannel),
						invokeRequestChannelHandler(),
					),
				)
				err = connectionApp.Err()
				if ctx != nil {
					err = multierr.Append(err, ctx.Err())
				}
				if err != nil {
					self.ZapLogger.Error(
						"Error in fxApp.Err() when creating NewConnectionInstance()",
						zap.Error(err))
					onErrorFlush(err)
					continue loop
				}

				if self.OnCreateConnection != nil {
					self.OnCreateConnection.OnCreateConnection(
						uniqueReference,
						connectionApp.Err(),
						ctx,
						cancelFunc)
				}

				// TODO: Adhere to timeouts
				err = connectionApp.Start(context.Background())
				if err != nil {
					onErrorFlush(err)
					continue loop
				}

				err = self.ConnectionReactor.AddAcceptedChannel(uniqueReference, connectionApp)
				if err != nil {
					// ??
				}
				// this function is part of the GoFunctionCounter count
				go func(
					app *fx.App,
					ctx context.Context,
					cancelFunc context.CancelFunc,
					connectionReactor internal.ISshConnectionReactor,
					uniqueReference string,
				) {
					functionName := self.GoFunctionCounter.CreateFunctionName("SshChannelListenerManager.ListenForNewConnections.WaitForConnection.Done")
					defer func(GoFunctionCounter GoFunctionCounter.IService, name string) {
						_ = GoFunctionCounter.Remove(name)
					}(self.GoFunctionCounter, functionName)
					_ = self.GoFunctionCounter.Add(functionName)

					//
					<-ctx.Done()
					var errList error
					errList = connectionReactor.RemoveAcceptedChannel(uniqueReference)
					// TODO: Adhere to timeouts
					errList = multierr.Append(errList, app.Stop(context.Background()))
					if errList != nil {
						self.ZapLogger.Error(
							"Stopping error. not really a problem. informational",
							zap.Error(errList))
					}
					if cancelFunc != nil {
						cancelFunc()
					}

				}(connectionApp, ctx, connCancelFunc, self.ConnectionReactor, uniqueReference)

				continue loop
			} else {
				err = acceptNewChannel.Reject(ssh.ResourceShortage, "")
				if err != nil {
					self.ZapLogger.Error("On acceptNewChannel.Reject", zap.Error(err))
				}
				continue loop
			}
		}
		self.ZapLogger.Info("Leaving accept loop")
	}
	var err error = nil
	// check if connection manager state is IFxService.Started
	actualState := self.ConnectionManager.State()
	if actualState != IFxService.Started {
		newError := IFxService.NewServiceStateError(
			self.ConnectionManager.ServiceName(),
			"Failed to start connection Listener",
			IFxService.Started,
			actualState)
		err = multierr.Append(err, newError)
	}

	if err == nil {
		// this function is part of the GoFunctionCounter count
		go func() {
			functionName := self.GoFunctionCounter.CreateFunctionName("SshChannelListenerManager.ListenForNewConnections.Accept")
			defer func(GoFunctionCounter GoFunctionCounter.IService, name string) {
				_ = GoFunctionCounter.Remove(name)
			}(self.GoFunctionCounter, functionName)
			_ = self.GoFunctionCounter.Add(functionName)

			//
			goFunc()
		}()
	}
	return err
}

func (self *manager) acceptWithContext() (common.ISshNewChannel, context.CancelFunc, error) {
	return self.listener.AcceptWithContext()
}

func NewManager(
	params struct {
		fx.In
		UseProxy                                 bool     `name:"UseProxy"`
		ConnectionUrl                            *url.URL `name:"ConnectionUrl"`
		ProxyUrl                                 *url.URL `name:"ProxyUrl"`
		ListenerAccept                           ISshListenerAccept
		OnCreateConnection                       goCommsDefinitions.IOnCreateConnection
		ConnectionManager                        goConnectionManager.IService
		CancelCtx                                context.Context
		CancelFunction                           context.CancelFunc
		Settings                                 *sshChannelListenerManagerSettings
		ZapLogger                                *zap.Logger
		StackName                                string `name:"StackName"`
		ConnectionName                           string `name:"ConnectionName"`
		ConnectionInstancePrefix                 string `name:"ConnectionInstancePrefix"`
		UniqueSessionNumber                      interfaces.IUniqueReferenceService
		AdditionalFxOptionsForConnectionInstance func() fx.Option
		ConnectionReactor                        internal.ISshConnectionReactor
		SshChannelSettings                       common.ISshChannelSettings `optional:"true"`
		GoFunctionCounter                        GoFunctionCounter.IService
	},
) (*manager, error) {

	if params.ConnectionManager.State() != IFxService.Started {
		return nil, IFxService.NewServiceStateError(
			params.ConnectionManager.ServiceName(),
			"Service in incorrect state", IFxService.Started,
			params.ConnectionManager.State())
	}

	netManager, err := netBase.NewNetManager(
		params.ConnectionName,
		params.ConnectionInstancePrefix,
		params.UseProxy,
		params.ProxyUrl,
		params.ConnectionUrl,
		params.CancelCtx,
		params.CancelFunction,
		params.StackName,
		params.ConnectionManager,
		params.Settings.userContext,
		params.ZapLogger,
		params.UniqueSessionNumber,
		params.AdditionalFxOptionsForConnectionInstance,
		params.GoFunctionCounter,
	)
	if err != nil {
		return nil, err
	}

	return &manager{
		ConnNetManager: netBase.ConnNetManager{
			NetManager: netManager,
		},
		listener:           params.ListenerAccept,
		MaxConnections:     params.Settings.MaxConnections,
		OnCreateConnection: params.OnCreateConnection,
		ConnectionReactor:  params.ConnectionReactor,
		SshChannelSettings: params.SshChannelSettings,
	}, nil
}
