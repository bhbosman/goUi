module github.com/bhbosman/goUi

go 1.18

require (
	github.com/bhbosman/gocommon v0.0.0-20220621055214-3b04298a9d45
	github.com/bhbosman/goerrors v0.0.0-20210201065523-bb3e832fa9ab // indirect
	github.com/cskr/pubsub v1.0.2
	github.com/gdamore/tcell/v2 v2.5.1
	github.com/golang/mock v1.4.4
	github.com/icza/gox v0.0.0-20220321141217-e2d488ab2fbc // indirect
	github.com/rivo/tview v0.0.0-20220307222120-9994674d60a8
	go.uber.org/fx v1.14.2
	go.uber.org/multierr v1.6.0
	go.uber.org/zap v1.21.0
)

require (
	github.com/gdamore/encoding v1.0.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/dig v1.12.0 // indirect
	golang.org/x/sys v0.0.0-20220318055525-2edf467146b5 // indirect
	golang.org/x/term v0.0.0-20210220032956-6a3ed077a48d // indirect
	golang.org/x/text v0.3.7 // indirect
)

replace github.com/gdamore/tcell/v2 => github.com/bhbosman/tcell/v2 v2.5.2-0.20220624055704-f9a9454fab5b

replace github.com/golang/mock => ../gomock


replace github.com/bhbosman/gocommon => ../gocommon
