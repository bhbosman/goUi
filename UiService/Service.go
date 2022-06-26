package UiService

import (
	"context"
	"fmt"
	"github.com/bhbosman/goUi/ui"
	"github.com/bhbosman/gocommon/Services/IFxService"
	"github.com/cskr/pubsub"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"sort"
	"strconv"
	"strings"
)

type Service struct {
	pubSub *pubsub.PubSub
	state  IFxService.State
}

func (self *Service) ServiceName() string {
	return "UiService"
}

func (self *Service) State() IFxService.State {
	return self.state
}

func (self *Service) OnStart(_ context.Context) error {
	self.state = IFxService.Started
	return nil
}

func (self *Service) OnStop(_ context.Context) error {
	//self.cancelFunc()
	self.state = IFxService.Stopped
	return nil
}

func (self *Service) Build(
	app *tview.Application,
	registeredSlides ...ui.ISlideFactory,
) (ui.IPrimitiveCloser, error) {
	return self.BuildApp(app, registeredSlides...)
}

func (self *Service) BuildApp(
	app *tview.Application,
	slideFactories ...ui.ISlideFactory,
) (ui.IPrimitiveCloser, error) {
	m := make(map[string]bool)
	for _, slide := range slideFactories {
		if _, ok := m[slide.Title()]; ok {
			return nil, fmt.Errorf("multiple slideFactories with name %v", slide.Title())
		}
		m[slide.Title()] = true
	}

	sort.Slice(
		slideFactories,
		func(i, j int) bool {
			iOrderNumber := slideFactories[i].OrderNumber()
			jOrderNumber := slideFactories[j].OrderNumber()
			b := iOrderNumber < jOrderNumber
			if !b {
				return false
			}
			b = iOrderNumber > jOrderNumber
			if !b {
				return true
			}
			return strings.Compare(slideFactories[i].Title(), slideFactories[j].Title()) == -1
		})
	pages := tview.NewPages()
	info := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false).
		SetHighlightedFunc(func(added, removed, remaining []string) {
			pages.SwitchToPage(added[0])
		})

	nextSlide := func() {
		slide, _ := strconv.Atoi(info.GetHighlights()[0])
		slide = (slide + 1) % len(slideFactories)
		info.Highlight(strconv.Itoa(slide)).
			ScrollToHighlight()
	}
	previousSlide := func() {
		slide, _ := strconv.Atoi(info.GetHighlights()[0])
		slide = (slide - 1 + len(slideFactories)) % len(slideFactories)
		info.Highlight(strconv.Itoa(slide)).
			ScrollToHighlight()
	}

	var closers []ui.IPrimitiveCloser
	for index, slide := range slideFactories {
		title, primitive, _ := slide.Content()(nextSlide)
		closers = append(closers, primitive)
		pages.AddPage(strconv.Itoa(index), primitive, true, index == 0)
		_, _ = fmt.Fprintf(info, `%d ["%d"][green]%s[white][""]  `, index+1, index, title)
	}
	info.Highlight("0")

	// Create the main layout.
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(pages, 0, 1, true).
		AddItem(info, 1, 1, false)

	app.SetInputCapture(
		func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyCtrlO {
				nextSlide()
				return nil
			} else if event.Key() == tcell.KeyCtrlP {
				previousSlide()
				return nil
			}
			return event
		})

	return ui.NewPrimitiveWithCloser(layout, closers), nil
}

func NewService(pubSub *pubsub.PubSub) *Service {
	result := &Service{
		pubSub: pubSub,
	}
	return result
}
