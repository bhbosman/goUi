package cmSlide

import (
	"github.com/bhbosman/goUi/ui"
	"github.com/rivo/tview"
)

type Factory struct {
	app     *tview.Application
	service *Service
}

func (self *Factory) OrderNumber() int {
	return 100
}

func NewFactory(
	app *tview.Application,
	service *Service,
) (*Factory, error) {
	return &Factory{
		app:     app,
		service: service,
	}, nil
}

func (self *Factory) Title() string {
	return "Connections"
}

func (self *Factory) Content(nextSlide func()) (string, ui.IPrimitiveCloser, error) {
	slide, err := NewConnectionSlide(
		self.app,
		self.service,
	)
	if err != nil {
		return "", nil, err
	}
	return self.Title(), slide, nil
}
