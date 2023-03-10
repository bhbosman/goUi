package cmSlide

import (
	"github.com/bhbosman/goConnectionManager"
	"github.com/bhbosman/goUi/ui"
	"github.com/rivo/tview"
)

type factory struct {
	app               *tview.Application
	service           *Service
	connectionManager goConnectionManager.IService
}

func (self *factory) OrderNumber() int {
	return 100
}

func (self *factory) Title() string {
	return "Connections"
}

func (self *factory) Content(nextSlide func()) (string, ui.IPrimitiveCloser, error) {
	slide, err := newConnectionSlide(
		self.app,
		self.service,
		self.connectionManager,
	)
	if err != nil {
		return "", nil, err
	}
	return self.Title(), slide, nil
}

func newFactory(
	app *tview.Application,
	service *Service,
	connectionManager goConnectionManager.IService,
) (*factory, error) {
	return &factory{
		app:               app,
		service:           service,
		connectionManager: connectionManager,
	}, nil
}
