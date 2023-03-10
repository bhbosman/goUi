package cmSlide

import (
	"go.uber.org/fx"
	"golang.org/x/net/context"
)

func InvokeConnectionManagerSlide() fx.Option {
	return fx.Options(
		fx.Invoke(
			func(
				params struct {
					fx.In
					Service   IConnectionSlideService
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
