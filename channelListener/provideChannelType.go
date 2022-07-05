package channelListener

import "go.uber.org/fx"

func provideChannelType(channelType string) fx.Option {
	return fx.Provide(
		fx.Annotated{
			Name: "ChannelType",
			Target: func() string {
				return channelType
			},
		},
	)
}
