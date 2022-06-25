package ConnectionSlide

import (
	"context"
	"github.com/bhbosman/goUi/ui"
	"github.com/cskr/pubsub"
	"github.com/rivo/tview"
)

type Factory struct {
	applicationContext context.Context
	pubSub             *pubsub.PubSub
	app                *tview.Application
}

func (self *Factory) OrderNumber() int {
	return 100
}

func NewFactory(
	applicationContext context.Context,
	pubSub *pubsub.PubSub,
	app *tview.Application,
) (*Factory, error) {
	return &Factory{
		applicationContext: applicationContext,
		pubSub:             pubSub,
		app:                app,
	}, nil
}

func (self *Factory) Title() string {
	return "Connections"
}

func (self *Factory) Content() ui.SlideCallback {
	return func(nextSlide func()) (string, ui.IPrimitiveCloser) {
		return self.Title(), NewConnectionSlide(self.applicationContext, self.pubSub, self.app)

	}
}
