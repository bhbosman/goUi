package GoFunctionCounterSlide

import (
	ui2 "github.com/bhbosman/goUi/ui"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/rivo/tview"
)

type CoverSlideFactory struct {
	app     *tview.Application
	service GoFunctionCounter.IService
}

func NewCoverSlideFactory(
	app *tview.Application,
	service GoFunctionCounter.IService,
) *CoverSlideFactory {
	return &CoverSlideFactory{
		app:     app,
		service: service,
	}
}

func (self *CoverSlideFactory) OrderNumber() int {
	return 2
}

func (self *CoverSlideFactory) Content(nextSlide func()) (string, ui2.IPrimitiveCloser, error) {
	primitive := NewSlide(
		self.app,
	)
	primitive.init()
	self.service.SetConnectionListChange(primitive.SetConnectionListChange)

	return self.Title(), primitive, nil
}

func (self *CoverSlideFactory) Title() string {
	return "Go Function Counter"
}
