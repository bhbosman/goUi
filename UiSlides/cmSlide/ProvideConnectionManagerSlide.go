package cmSlide

import (
	"context"
	"github.com/bhbosman/goConnectionManager"
	"github.com/bhbosman/goUi/ui"
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
					UniqueReferenceService  interfaces.IUniqueReferenceService
					Logger                  *zap.Logger
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
						App     *tview.Application
						Service *Service
					},
				) (ui.ISlideFactory, error) {

					return NewFactory(
						params.App,
						params.Service,
					)
				},
			},
		),
	)
}
