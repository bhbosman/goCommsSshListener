module github.com/bhbosman/goCommsSshListener

go 1.18

require (
	github.com/bhbosman/goCommsDefinitions v0.0.0-20230329100608-a6a24c060ad8
	github.com/bhbosman/goCommsNetListener v0.0.0-20220611182354-46d10d89b8e1
	github.com/bhbosman/goCommsStacks v0.0.0-20220611141421-a7d405cadbfa
	github.com/bhbosman/goConnectionManager v0.0.0-20230329104211-b2d06385b410
	github.com/bhbosman/gocommon v0.0.0-20230329101749-40db0f52d859
	github.com/bhbosman/gocomms v0.0.0-20230329110556-946ebc6ff5f4
	github.com/bhbosman/goerrors v0.0.0-20220623084908-4d7bbcd178cf
	github.com/bhbosman/gomessageblock v0.0.0-20230308173223-e8144f25444c
	github.com/bhbosman/goprotoextra v0.0.2
	github.com/gdamore/tcell/v2 v2.5.1
	github.com/golang/mock v1.6.0
	go.uber.org/fx v1.19.2
	go.uber.org/multierr v1.10.0
	go.uber.org/zap v1.24.0
	golang.org/x/crypto v0.7.0
	golang.org/x/net v0.8.0
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
)

require (
	github.com/cskr/pubsub v1.0.2 // indirect
	golang.org/x/term v0.6.0 // indirect
)

require (
	github.com/cenkalti/backoff/v4 v4.1.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emirpasic/gods v1.12.0 // indirect
	github.com/icza/gox v0.0.0-20220321141217-e2d488ab2fbc // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/reactivex/rxgo/v2 v2.5.0
	github.com/stretchr/objx v0.4.0 // indirect
	github.com/stretchr/testify v1.8.0 // indirect
	github.com/teivah/onecontext v0.0.0-20200513185103-40f981bfd775 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/dig v1.16.1 // indirect
	golang.org/x/sys v0.6.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/gdamore/tcell/v2 => github.com/bhbosman/tcell/v2 v2.5.2-0.20220624055704-f9a9454fab5b

replace github.com/golang/mock => github.com/bhbosman/gomock v1.6.1-0.20230302060806-d02c40b7514e

replace github.com/cskr/pubsub => github.com/bhbosman/pubsub v1.0.3-0.20220802200819-029949e8a8af

replace github.com/rivo/tview => github.com/bhbosman/tview v0.0.0-20230310100135-f8b257a85d36

replace github.com/bhbosman/gocomms => ../gocomms

replace github.com/bhbosman/goMessages => ../goMessages

replace github.com/bhbosman/goCommsDefinitions => ../goCommsDefinitions

replace github.com/bhbosman/goCommsSSHProtocols => ../goCommsSSHProtocols

replace github.com/bhbosman/goCommsStacks => ../goCommsStacks

replace github.com/bhbosman/goCommsNetListener => ../goCommsNetListener

replace github.com/bhbosman/goConnectionManager => ../goConnectionManager

replace github.com/bhbosman/goprotoextra => ../goprotoextra
