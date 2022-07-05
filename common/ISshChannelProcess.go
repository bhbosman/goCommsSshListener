package common

import (
	"github.com/bhbosman/goprotoextra"
	"github.com/gdamore/tcell/v2/terminfo"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"io"
)

type ISetSize interface {
	SetSize(cols int, rows int) error
}

type ISshChannelProcess interface {
	io.WriteCloser
	ISetSize
	RunHandler() error
	SetLookupTerminfo(terminfo *terminfo.Terminfo) error
}

type SshBuildInChannelProcess int

const (
	NoneInChannelProcess SshBuildInChannelProcess = iota
	UseCustomChannelProcess
	EchoBuildInChannelProcess
)

type ISshChannelSettings interface {
	GetSettingsFor(sessionType string) (ISettings, error)
}

type ISettings interface {
	Name() string
}

type ISshChannelSessionSettings interface {
	ISettings
	CreateChannelProcess(
		sshChannel ISshChannel,
		parentCtx context.Context,
		parentCancelFunc context.CancelFunc,
		onSend goprotoextra.ToConnectionFunc,
		onSendReplacement rxgo.NextFunc,
		logger *zap.Logger,
	) (ISshChannelProcess, error)
	UseDefault() (SshBuildInChannelProcess, error)
}
