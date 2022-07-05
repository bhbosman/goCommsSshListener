package channelListener

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goCommsStacks/bottom"
	"github.com/bhbosman/goCommsStacks/top"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/common"
	"go.uber.org/fx"
	"golang.org/x/crypto/ssh"
	"net/url"
	"time"
)

func NewManagerApp(
	name string,
	serviceIdentifier model.ServiceIdentifier,
	serviceDependentOn model.ServiceIdentifier,
	connectionInstancePrefix string,
	urlAsString string,
	channels <-chan ssh.NewChannel,
	conn goCommsDefinitions.ISpecificInformationForConnection,
	cancelFunc context.CancelFunc,
	settings ...common.INetManagerSettingsApply,
) common.NetAppFuncInParamsCallback {
	return func(params common.NetAppFuncInParams) messages.CreateAppCallback {
		return messages.CreateAppCallback{
			ServiceId:         serviceIdentifier,
			ServiceDependency: serviceDependentOn,
			Name:              name,
			Callback: func() (*fx.App, context.CancelFunc, error) {
				resultCancelFunc := func() {
					cancelFunc()
				}
				netListenSettings := &sshChannelListenerManagerSettings{
					NetManagerSettings:    common.NewNetManagerSettings(512),
					userContext:           nil,
					netListenerFactory:    provideCreateListenResource,
					listenerAcceptFactory: provideCreateListenAcceptResource,
				}
				netListenSettings.AddFxOptionsForConnectionInstance(
					[]fx.Option{
						goCommsDefinitions.ProvideTransportFactoryForEmptyName(
							top.ProvideTopStack(),
							bottom.Provide(),
						),
						goCommsDefinitions.ProvideTransportFactoryForSshChannelSession(
							top.ProvideTopStack(),
							//session.ProvideSshSessionProtocolStack(),
							bottom.Provide(),
						),
					},
				)

				for _, setting := range settings {
					if setting == nil {
						continue
					}
					if listenAppSettingsApply, ok := setting.(ISshListenAppSettingsApply); ok {
						err := listenAppSettingsApply.apply(netListenSettings)
						if err != nil {
							return nil, resultCancelFunc, err
						}
					} else {
						err := setting.ApplyNetManagerSettings(&netListenSettings.NetManagerSettings)
						if err != nil {
							return nil, resultCancelFunc, err
						}
					}
				}

				callbackForConnectionInstance, err := netListenSettings.Build()
				if err != nil {
					return nil, nil, err
				}

				sshUrl, err := url.Parse(urlAsString)
				if err != nil {
					return nil, nil, err
				}
				options := common.ConnectionApp(
					time.Hour,
					time.Hour,
					name,
					connectionInstancePrefix,
					false,
					nil,
					sshUrl,
					goCommsDefinitions.TransportFactoryEmptyName,
					params,
					callbackForConnectionInstance,
					fx.Options(netListenSettings.MoreOptions...),
					fx.Supply(netListenSettings),
					fx.Provide(fx.Annotated{Target: NewManager}),
					fx.Provide(fx.Annotated{Target: netListenSettings.OnCreateConnectionFactory}),
					fx.Provide(fx.Annotated{Target: netListenSettings.listenerAcceptFactory}),
					fx.Provide(fx.Annotated{Target: netListenSettings.netListenerFactory}),
					fx.Provide(
						fx.Annotated{
							Target: func() (goCommsDefinitions.ISpecificInformationForConnection, error) {
								return conn, nil
							},
						},
					),
					fx.Provide(
						fx.Annotated{
							Target: func() (<-chan ssh.NewChannel, error) {
								return channels, nil
							},
						}),
					fx.Invoke(sshInvokeListenForNewConnections),
				)
				fxApp := fx.New(options)
				return fxApp, resultCancelFunc, fxApp.Err()
			},
		}
	}
}
