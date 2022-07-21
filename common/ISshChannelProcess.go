package common

import (
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/gdamore/tcell/v2/terminfo"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"io"
)

type ISetSize interface {
	SetSize(cols int, rows int) error
}

type IChannelProcess interface {
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
		sshChannel IChannel,
		parentCtx context.Context,
		parentCancelFunc context.CancelFunc,
		onSend rxgo.NextFunc,
		logger *zap.Logger,
		GoFunctionCounter GoFunctionCounter.IService,
	) (IChannelProcess, error)
	UseDefault() (SshBuildInChannelProcess, error)
}
