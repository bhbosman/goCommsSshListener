package channelListener

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goCommsStacks/bottom"
	"github.com/bhbosman/goCommsStacks/topStack"
	"github.com/bhbosman/gocommon"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocomms/common"
	"go.uber.org/fx"
	"golang.org/x/crypto/ssh"
	"net/url"
	"time"
)

func NewManagerApp(
	name string,
	connectionInstancePrefix string,
	urlAsString string,
	channels <-chan ssh.NewChannel,
	conn goCommsDefinitions.ISpecificInformationForConnection,
	settings ...common.INetManagerSettingsApply,
) common.NetAppFuncInParamsCallback {
	return func(params common.NetAppFuncInParams) messages.CreateAppCallback {
		return messages.CreateAppCallback{
			Name: name,
			Callback: func() (messages.IApp, gocommon.ICancellationContext, error) {
				netListenSettings := &channelListenerManagerSettings{
					NetManagerSettings:    common.NewNetManagerSettings(512),
					userContext:           nil,
					netListenerFactory:    provideCreateListenResource,
					listenerAcceptFactory: provideCreateListenAcceptResource,
				}
				netListenSettings.AddFxOptionsForConnectionInstance(
					[]fx.Option{
						goCommsDefinitions.ProvideTransportFactoryForEmptyName(
							topStack.Provide(),
							bottom.Provide(),
						),
						goCommsDefinitions.ProvideTransportFactoryForSshChannelSession(
							topStack.Provide(),
							bottom.Provide(),
						),
					},
				)

				namedLogger := params.ZapLogger.Named(name)
				ctx, cancelFunc := context.WithCancel(params.ParentContext)
				cancellationContext, err := gocommon.NewCancellationContextNoCloser(name, cancelFunc, ctx, namedLogger)

				for _, setting := range settings {
					if setting == nil {
						continue
					}
					if listenAppSettingsApply, ok := setting.(iListenAppSettingsApply); ok {
						err := listenAppSettingsApply.apply(netListenSettings)
						if err != nil {
							return nil, cancellationContext, err
						}
					} else {
						err := setting.ApplyNetManagerSettings(&netListenSettings.NetManagerSettings)
						if err != nil {
							return nil, cancellationContext, err
						}
					}
				}

				callbackForConnectionInstance, err := netListenSettings.Build()
				if err != nil {
					return nil, nil, err
				}

				var sshUrl *url.URL
				sshUrl, err = url.Parse(urlAsString)
				if err != nil {
					return nil, nil, err
				}

				options := common.ConnectionApp(
					time.Hour,
					time.Hour,
					name,
					connectionInstancePrefix,
					params,
					cancellationContext,
					namedLogger,
					callbackForConnectionInstance,
					fx.Options(netListenSettings.MoreOptions...),
					fx.Supply(netListenSettings),
					goCommsDefinitions.ProvideUrl("ConnectionUrl", sshUrl),
					goCommsDefinitions.ProvideUrl("ProxyUrl", nil),
					goCommsDefinitions.ProvideBool("UseProxy", false),
					fx.Provide(fx.Annotated{Target: NewManager}),
					fx.Provide(fx.Annotated{Target: netListenSettings.listenerAcceptFactory}),
					fx.Provide(fx.Annotated{Target: netListenSettings.netListenerFactory}),
					provideISpecificInformationForConnection(conn),
					provideNewChannelChannel(channels),
					invokeListenForNewConnections(),
				)
				fxApp := fx.New(options)
				return fxApp, cancellationContext, fxApp.Err()
			},
		}
	}
}
