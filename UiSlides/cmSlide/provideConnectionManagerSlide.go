package cmSlide

import (
	"context"
	"github.com/bhbosman/goConnectionManager"
	"github.com/bhbosman/goUi/ui"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/services/interfaces"
	"github.com/cskr/pubsub"
	"github.com/rivo/tview"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ConnectionManagerSlideResult struct {
	fx.Out
	PrimitiveCloserForGroup ui.IPrimitiveCloser `group:"RegisteredMainWindowSlides"`
	PrimitiveCloserForName  ui.IPrimitiveCloser `name:"ConnectionSlide"`
	Slide                   *slide
}

func ProvideConnectionManagerSlide() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Target: func(
					params struct {
						fx.In
						PubSub                  *pubsub.PubSub  `name:"Application"`
						ApplicationContext      context.Context `name:"Application"`
						ConnectionManagerHelper goConnectionManager.IHelper
						ConnectionManager       goConnectionManager.IService
						UniqueReferenceService  interfaces.IUniqueReferenceService
						Logger                  *zap.Logger
						GoFunctionCounter       GoFunctionCounter.IService
					},
				) (IConnectionSlideService, error) {
					s, e := newService(
						params.ApplicationContext,
						params.PubSub,
						func() (IConnectionSlideData, error) {
							return NewData()
						},
						params.ConnectionManagerHelper,
						params.UniqueReferenceService,
						params.Logger,
						params.GoFunctionCounter,
						params.ConnectionManager,
					)
					if e != nil {
						return nil, e
					}
					return s, nil
				},
			},
		),
		fx.Provide(
			fx.Annotated{
				Target: func(
					params struct {
						fx.In
						Lifecycle fx.Lifecycle
						App       *tview.Application
						Service   IConnectionSlideService
					},
				) (ConnectionManagerSlideResult, error) {
					connectionSlide, err := newConnectionSlide(
						"Connections",
						params.App,
						params.Service,
					)
					if err != nil {
						return ConnectionManagerSlideResult{}, err
					}
					return ConnectionManagerSlideResult{
						PrimitiveCloserForGroup: connectionSlide,
						PrimitiveCloserForName:  connectionSlide,
						Slide:                   connectionSlide,
					}, nil
				},
			},
		),
	)
}
