package GoFunctionCounterSlide

import (
	ui2 "github.com/bhbosman/goUi/ui"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/rivo/tview"
)

type factory struct {
	app     *tview.Application
	service GoFunctionCounter.IService
}

func (self *factory) OrderNumber() int {
	return 2
}

func (self *factory) Content(nextSlide func()) (string, ui2.IPrimitiveCloser, error) {
	primitive := newSlide(
		self.app,
	)
	primitive.init()
	self.service.SetConnectionListChange(primitive.SetConnectionListChange)

	return self.Title(), primitive, nil
}

func (self *factory) Title() string {
	return "Go Function Counter"
}

func newGoFunctionSideFactory(
	app *tview.Application,
	service GoFunctionCounter.IService,
) *factory {
	return &factory{
		app:     app,
		service: service,
	}
}
