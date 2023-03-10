package cmSlide

import (
	"context"
	"github.com/bhbosman/goConnectionManager"
	"github.com/bhbosman/goUi/ui"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/Services/interfaces"
	"github.com/cskr/pubsub"
	"github.com/rivo/tview"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func InvokeConnectionManagerSlide() fx.Option {
	return fx.Options(
		fx.Invoke(
			func(
				params struct {
					fx.In
					Service   *Service
					Lifecycle fx.Lifecycle
				}) error {
				params.Lifecycle.Append(
					fx.Hook{
						OnStart: func(ctx context.Context) error {
							return params.Service.OnStart(ctx)
						},
						OnStop: func(ctx context.Context) error {
							return params.Service.OnStop(ctx)
						},
					})
				return nil
			},
		),
	)
}

func ProvideConnectionManagerSlide() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Target: func(params struct {
					fx.In
					PubSub                  *pubsub.PubSub  `name:"Application"`
					ApplicationContext      context.Context `name:"Application"`
					ConnectionManagerHelper goConnectionManager.IHelper
					ConnectionManager       goConnectionManager.IService
					UniqueReferenceService  interfaces.IUniqueReferenceService
					Logger                  *zap.Logger
					GoFunctionCounter       GoFunctionCounter.IService
				}) (*Service, error) {
					s, e := NewService(
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
				Group: "RegisteredMainWindowSlides",
				Target: func(
					params struct {
						fx.In
						App               *tview.Application
						Service           *Service
						ConnectionManager goConnectionManager.IService
					},
				) (ui.ISlideFactory, error) {

					return newFactory(
						params.App,
						params.Service,
						params.ConnectionManager,
					)
				},
			},
		),
	)
}
