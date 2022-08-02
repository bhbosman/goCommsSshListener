package channelListener

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	"go.uber.org/fx"
	"golang.org/x/crypto/ssh"
	url2 "net/url"
)

func provideCreateListenResource(
	params struct {
		fx.In
		Lifecycle                        fx.Lifecycle
		UseProxy                         bool      `name:"UseProxy"`
		ConnectionUrl                    *url2.URL `name:"ConnectionUrl"`
		ProxyUrl                         *url2.URL `name:"ProxyUrl"`
		CancelContext                    context.Context
		Channel                          <-chan ssh.NewChannel
		SpecificInformationForConnection goCommsDefinitions.ISpecificInformationForConnection
	},
) (iChannelListener, error) {
	con, err := newChannelListener(
		params.CancelContext,
		params.Channel,
		params.SpecificInformationForConnection,
		params.ConnectionUrl)
	if err != nil {
		return nil, err
	}
	params.Lifecycle.Append(fx.Hook{
		OnStart: nil,
		OnStop: func(ctx context.Context) error {
			return con.Close()
		},
	})
	return con, nil
}
