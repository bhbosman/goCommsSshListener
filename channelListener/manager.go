package channelListener

import (
	"context"
	"github.com/bhbosman/goCommsSshListener/common"
	"github.com/bhbosman/goCommsSshListener/internal"
	"github.com/bhbosman/goConnectionManager"
	"github.com/bhbosman/gocommon"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/services/IFxService"
	"github.com/bhbosman/gocommon/services/interfaces"
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
	listener           iListenerAccept
	MaxConnections     int
	ConnectionReactor  internal.ISshConnectionReactor
	SshChannelSettings common.ISshChannelSettings
}

func (self *manager) ListenForNewConnections() error {
	actualState := self.ConnectionManager.State()
	if actualState != IFxService.Started {
		newError := IFxService.NewServiceStateError(
			self.ConnectionManager.ServiceName(),
			"Failed to start connection Listener",
			IFxService.Started,
			actualState)
		return newError
	}
	return self.GoFunctionCounter.GoRun(
		"SshChannelListenerManager.ListenForNewConnections.Accept",
		func() {
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
					break loop
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
					var acceptedChannel common.IChannel
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

					acceptedChannel = newSshChannelWithSemaphoreWrapper(acceptedChannel, sem)

					uniqueReference := self.UniqueSessionNumber.Next(self.ConnectionInstancePrefix)

					ctx, cancelFunc := context.WithCancel(self.CancellationContext.CancelContext())
					cancellationContext, err := gocommon.NewCancellationContext(
						uniqueReference,
						cancelFunc,
						ctx,
						self.ZapLogger,
						acceptedChannel)
					if err != nil {
						return
					}

					connectionInstance := netBase.NewConnectionInstance(
						self.ConnectionUrl,
						self.UniqueSessionNumber,
						self.ConnectionManager,
						cancellationContext,
						self.AdditionalFxOptionsForConnectionInstance,
						self.ZapLogger,
					)
					connectionApp, err := connectionInstance.NewConnectionInstanceWithStackName(
						uniqueReference,
						self.GoFunctionCounter,
						model.ServerConnection,
						acceptedChannel,
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
									) (common.IChannel, error) {
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
					if err != nil {
						cancellationContext.CancelWithError("sadsadasd", err)
					}
					onErr := func() {
						if connCancelFunc != nil {
							connCancelFunc()
						}
						if cancellationContext != nil {
							cancellationContext.Cancel("123")
						}
						err = multierr.Append(err, acceptedChannel.Close())
					}
					if err != nil {
						onErr()
						return
					}

					err = connectionApp.Start(context.Background())
					if err != nil {
						onErr()
						continue loop
					}

					err = self.ConnectionReactor.AddAcceptedChannel(uniqueReference, connectionApp)
					if err != nil {
						// ??
					}
					_ = gocommon.RegisterConnectionShutdown(
						uniqueReference,
						func(
							connectionReactor internal.ISshConnectionReactor,
							zapLogger *zap.Logger,
						) func() {
							return func() {
								var errList error
								errList = connectionReactor.RemoveAcceptedChannel(uniqueReference)
								// TODO: Adhere to timeouts
								errList = multierr.Append(errList, connectionApp.Stop(context.Background()))
								if errList != nil {
									zapLogger.Error(
										"Stopping error. not really a problem. informational",
										zap.Error(errList))
								}
								if connCancelFunc != nil {
									connCancelFunc()
								}
								err := multierr.Append(err, acceptedChannel.Close())
								err = multierr.Append(err,
									self.GoFunctionCounter.GoRun(
										"SshChannelListenerManager.ListenForNewConnections.Flush.AcceptedChannelRequestChannel",
										func() {
											for range acceptedChannelRequestChannel {
											}
										},
									),
								)
								self.ZapLogger.Error("OnErrorFlush", zap.Error(err))
							}
						}(
							self.ConnectionReactor,
							self.ZapLogger,
						),
						cancellationContext,
						self.CancellationContext,
					)
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
		},
	)
}

func (self *manager) acceptWithContext() (common.INewChannel, context.CancelFunc, error) {
	return self.listener.acceptWithContext()
}

func NewManager(
	params struct {
		fx.In
		UseProxy                                 bool     `name:"UseProxy"`
		ConnectionUrl                            *url.URL `name:"ConnectionUrl"`
		ProxyUrl                                 *url.URL `name:"ProxyUrl"`
		ListenerAccept                           iListenerAccept
		ConnectionManager                        goConnectionManager.IService
		CancelCtx                                context.Context
		CancelFunction                           context.CancelFunc
		Settings                                 *channelListenerManagerSettings
		ZapLogger                                *zap.Logger
		CancellationContext                      gocommon.ICancellationContext
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
		params.CancellationContext,
		params.ConnectionManager,
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
		ConnectionReactor:  params.ConnectionReactor,
		SshChannelSettings: params.SshChannelSettings,
	}, nil
}
