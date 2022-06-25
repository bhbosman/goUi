package UiService

import (
	"context"
	"github.com/bhbosman/goUi/UiSlides/connectionSlide"
	"github.com/bhbosman/goUi/UiSlides/intoductionSlide"
	"github.com/bhbosman/goUi/ui"
	"github.com/cskr/pubsub"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"go.uber.org/fx"
)

func ProvideTerminalApplication() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Target: func(
					params struct {
						fx.In
						ApplicationContext context.Context `name:"Application"`
						PubSub             *pubsub.PubSub  `name:"Application"`
					},
				) IUiService {
					return NewService(params.ApplicationContext, params.PubSub)
				},
			}),
		fx.Provide(
			fx.Annotated{
				Group: "RegisteredMainWindowSlides",
				Target: func(
					params struct {
						fx.In
					},
				) (ui.ISlideFactory, error) {
					return &intoductionSlide.CoverSlide{}, nil
				}}),
		fx.Provide(
			fx.Annotated{
				Group: "RegisteredMainWindowSlides",
				Target: func(
					params struct {
						fx.In
						App                *tview.Application
						ApplicationContext context.Context `name:"Application"`
						PubSub             *pubsub.PubSub  `name:"Application"`
					},
				) (ui.ISlideFactory, error) {
					return ConnectionSlide.NewFactory(params.ApplicationContext, params.PubSub, params.App)
				}}),
		fx.Provide(
			fx.Annotated{
				Target: func(
					params struct {
						fx.In
						UiApp                      IUiService
						App                        *tview.Application
						Screen                     tcell.Screen       `optional:"true"`
						RegisteredMainWindowSlides []ui.ISlideFactory `group:"RegisteredMainWindowSlides"`
					},
				) (ui.IPrimitiveCloser, error) {
					return params.UiApp.Build()(params.App, params.RegisteredMainWindowSlides...)

				},
			}),
		fx.Provide(
			fx.Annotated{
				Target: func(
					params struct {
						fx.In
					},
				) (*tview.Application, error) {
					result := tview.NewApplication()
					return result, nil
				},
			}),
	)
}
