package channelListener

import (
	"context"
	"github.com/bhbosman/goCommsSshListener/common"
)

type ISshListenerAccept interface {
	AcceptWithContext() (common.ISshNewChannel, context.CancelFunc, error)
}
