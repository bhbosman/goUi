package UiService

import (
	"github.com/bhbosman/goUi/ui"
	"github.com/cskr/pubsub"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"go.uber.org/fx"
)

func ProvideTerminalApplication__() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Target: func(
					params struct {
						fx.In
						PubSub *pubsub.PubSub `name:"Application"`
					},
				) IUiService {
					return NewService(params.PubSub)
				},
			},
		),

		fx.Provide(
			fx.Annotated{
				Target: func(
					params struct {
						fx.In
						UiApp                      IUiService
						App                        *tview.Application
						RegisteredMainWindowSlides []ui.ISlideFactory `group:"RegisteredMainWindowSlides"`
					},
				) (ui.IPrimitiveCloser, error) {
					return params.UiApp.Build(params.App, params.RegisteredMainWindowSlides...)
				},
			},
		),
		fx.Provide(
			fx.Annotated{
				Target: func(
					params struct {
						fx.In
						Screen tcell.Screen `optional:"true"`
					},
				) (*tview.Application, error) {
					result := tview.NewApplication()
					if params.Screen != nil {
						result = result.SetScreen(params.Screen)
					}
					return result, nil
				},
			},
		),
	)
}
