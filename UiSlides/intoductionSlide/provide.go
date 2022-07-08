package intoductionSlide

import (
	ui2 "github.com/bhbosman/goUi/ui"
	"go.uber.org/fx"
)

func ProvideCoverSlide() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Group: "RegisteredMainWindowSlides",
			Target: func(
				params struct {
					fx.In
				},
			) (ui2.ISlideFactory, error) {
				return &CoverSlideFactory{}, nil
			},
		},
	)
}
