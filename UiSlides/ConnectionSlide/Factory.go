package ConnectionSlide

import (
	"context"
	"github.com/bhbosman/goUi/ui"
	"github.com/bhbosman/gocommon/Services/IConnectionManager"
	"github.com/bhbosman/gocommon/Services/interfaces"
	"github.com/cskr/pubsub"
	"github.com/rivo/tview"
)

type Factory struct {
	applicationContext      context.Context
	pubSub                  *pubsub.PubSub
	app                     *tview.Application
	ConnectionManagerHelper IConnectionManager.IHelper
	UniqueReferenceService  interfaces.IUniqueReferenceService
}

func (self *Factory) OrderNumber() int {
	return 100
}

func NewFactory(
	applicationContext context.Context,
	pubSub *pubsub.PubSub,
	app *tview.Application,
	ConnectionManagerHelper IConnectionManager.IHelper,
	UniqueReferenceService interfaces.IUniqueReferenceService,
) (*Factory, error) {
	return &Factory{
		applicationContext:      applicationContext,
		pubSub:                  pubSub,
		app:                     app,
		ConnectionManagerHelper: ConnectionManagerHelper,
		UniqueReferenceService:  UniqueReferenceService,
	}, nil
}

func (self *Factory) Title() string {
	return "Connections"
}

func (self *Factory) Content() ui.SlideCallback {
	return func(nextSlide func()) (string, ui.IPrimitiveCloser, error) {
		slide, err := NewConnectionSlide(
			self.applicationContext,
			self.pubSub,
			self.app,
			self.ConnectionManagerHelper,
			self.UniqueReferenceService,
		)
		if err != nil {
			return "", nil, err
		}
		return self.Title(), slide, nil

	}
}
