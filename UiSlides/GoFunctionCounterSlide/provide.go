package GoFunctionCounterSlide

import (
	ui2 "github.com/bhbosman/goUi/ui"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/rivo/tview"
	"go.uber.org/fx"
)

func Provide() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Group: "RegisteredMainWindowSlides",
			Target: func(
				params struct {
					fx.In
					App     *tview.Application
					Service GoFunctionCounter.IService
				},
			) (ui2.ISlideFactory, error) {
				return newGoFunctionSideFactory(
					params.App,
					params.Service,
				), nil
			},
		},
	)
}
