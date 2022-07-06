package channelListener

import (
	"context"
	"github.com/bhbosman/goCommsSshListener/common"
)

type iListenerAccept interface {
	acceptWithContext() (common.INewChannel, context.CancelFunc, error)
}
