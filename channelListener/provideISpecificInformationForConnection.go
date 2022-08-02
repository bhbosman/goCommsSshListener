package channelListener

import (
	"github.com/bhbosman/goCommsDefinitions"
	"go.uber.org/fx"
)

func provideISpecificInformationForConnection(conn goCommsDefinitions.ISpecificInformationForConnection) fx.Option {
	return fx.Provide(
		fx.Annotated{
			Target: func() (goCommsDefinitions.ISpecificInformationForConnection, error) {
				return conn, nil
			},
		},
	)
}
