package netListener

import (
	"fmt"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goCommsNetListener"
	common2 "github.com/bhbosman/goCommsSshListener/common"
	internal2 "github.com/bhbosman/goCommsSshListener/internal"
	sshStack "github.com/bhbosman/goCommsSshListener/stack"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gocomms/intf"
	"go.uber.org/fx"
	"net/url"
)

func NewSshListenApp(
	name string,
	serviceIdentifier model.ServiceIdentifier,
	serviceDependentOn model.ServiceIdentifier,
	connectionInstancePrefix string,
	useProxy bool,
	proxyUrl *url.URL,
	connectionUrl *url.URL,
	channelSettings common2.ISshChannelSettings,
	settings ...common.INetManagerSettingsApply,
) common.NetAppFuncInParamsCallback {

	settings = append(
		settings,
		common.NewConnectionInstanceOptions(
			goCommsDefinitions.ProvideTransportFactoryForOnlySSHStack(
				sshStack.Provide(),
			),
			ProvideConnectionReactorFactory(),
			fx.Provide(
				fx.Annotated{
					Target: func() (common2.ISshChannelSettings, error) {
						return channelSettings, nil
					},
				},
			),
			fx.Provide(
				fx.Annotated{
					Target: func(
						params struct {
							fx.In
							ConnectionReactor intf.IConnectionReactor
						},
					) (internal2.ISshConnectionReactor, error) {
						if reactor, ok := params.ConnectionReactor.(internal2.ISshConnectionReactor); ok {
							return reactor, nil
						}
						return nil, fmt.Errorf("could not extract ISshConnectionReactor")
					},
				},
			),
		),
	)

	f := goCommsNetListener.NewNetListenApp(
		name,
		serviceIdentifier,
		serviceDependentOn,
		connectionInstancePrefix,
		useProxy,
		proxyUrl,
		connectionUrl,
		goCommsDefinitions.TransportFactoryOnlySSHStack,
		settings...,
	)
	return f
}
