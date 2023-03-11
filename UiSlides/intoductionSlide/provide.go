package intoductionSlide

import (
	"fmt"
	ui2 "github.com/bhbosman/goUi/ui"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"go.uber.org/fx"
	"strings"
)

func ProvideCoverSlide() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Group: "RegisteredMainWindowSlides",
			Target: func(
				params struct {
					fx.In
				},
			) (ui2.IPrimitiveCloser, error) {
				lines := strings.Split(logo, "\n")
				logoWidth := 0
				logoHeight := len(lines)
				for _, line := range lines {
					if len(line) > logoWidth {
						logoWidth = len(line)
					}
				}
				logoBox := tview.NewTextView().
					SetTextColor(tcell.ColorGreen)

				fmt.Fprint(logoBox, logo)

				// Create a frame for the subtitle and navigation infos.
				frame := tview.NewFrame(tview.NewBox()).
					SetBorders(0, 0, 0, 0, 0, 0).
					AddText(subtitle, true, tview.AlignCenter, tcell.ColorWhite).
					AddText("", true, tview.AlignCenter, tcell.ColorWhite).
					AddText(navigation, true, tview.AlignCenter, tcell.ColorDarkMagenta).
					AddText(mouse, true, tview.AlignCenter, tcell.ColorDarkMagenta)

				// Create a Flex layout that centers the logo and subtitle.
				flex := tview.NewFlex().
					SetDirection(tview.FlexRow).
					AddItem(tview.NewBox(), 0, 7, false).
					AddItem(tview.NewFlex().
						AddItem(tview.NewBox(), 0, 1, false).
						AddItem(logoBox, logoWidth, 1, true).
						AddItem(tview.NewBox(), 0, 1, false), logoHeight, 1, true).
					AddItem(frame, 0, 10, false)

				return ui2.NewPrimitiveNoCloser(-1, "Main", flex), nil
			},
		},
	)
}
