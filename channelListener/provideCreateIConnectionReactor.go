package channelListener

import (
	"context"
	"fmt"
	"github.com/bhbosman/goCommsSshListener/common"
	"github.com/bhbosman/goCommsSshListener/internal"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/goerrors"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func provideCreateIConnectionReactor() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Target: func(
				params struct {
					fx.In
					CancelCtx            context.Context
					CancelFunc           context.CancelFunc
					Logger               *zap.Logger
					SshConnectionReactor internal.ISshConnectionReactor
					SshChannelSettings   common.ISshChannelSettings `optional:"true"`
					SshChannel           common.IChannel
					GoFunctionCounter    GoFunctionCounter.IService
					ExtraData            []byte `name:"ChannelExtraData"`
					//ClientContext        interface{} `name:"UserContext"`
					ChannelType string `name:"ChannelType"`
				},
			) (intf.IConnectionReactor, error) {
				params.Logger.Info("Creating Connection reactor")
				switch params.ChannelType {
				case "session":
					if params.SshChannelSettings == nil {
						return nil, fmt.Errorf("settings required of type ISshChannelSettings")
					}
					settingsForChannelType, err := params.SshChannelSettings.GetSettingsFor(params.ChannelType)
					if err != nil {
						return nil, err
					}
					unk, ok := settingsForChannelType.(common.ISshChannelSessionSettings)
					if !ok {
						return nil, fmt.Errorf("settings for %v in correct. Must implement ISshChannelSessionSettings", params.ChannelType)
					}
					return NewReactor(
						params.ChannelType,
						params.CancelCtx,
						params.CancelFunc,
						params.Logger,
						params.SshChannel,
						unk,
						params.ExtraData,
						params.GoFunctionCounter,
					)
				case "direct-tcpip":
					return nil, goerrors.NewInvalidParamError(
						"channelType",
						fmt.Sprintf("ssh channel not supported %v", params.ChannelType))
				default:
					return nil, goerrors.NewInvalidParamError(
						"channelType",
						fmt.Sprintf("ssh channel not supported %v", params.ChannelType))
				}
			},
		},
	)
}
