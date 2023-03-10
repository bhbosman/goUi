package cmSlide

import (
	"github.com/bhbosman/goUi/ui"
	"github.com/rivo/tview"
)

type factory struct {
	app     *tview.Application
	service IConnectionSlideService
}

func (self *factory) OrderNumber() int {
	return 100
}

func (self *factory) Title() string {
	return "Connections"
}

func (self *factory) Content() (string, ui.IPrimitiveCloser, error) {
	slide, err := newConnectionSlide(
		self.app,
		self.service,
	)
	if err != nil {
		return "", nil, err
	}
	return self.Title(), slide, nil
}

func newFactory(
	app *tview.Application,
	service IConnectionSlideService,
) (*factory, error) {
	return &factory{
		app:     app,
		service: service,
	}, nil
}
