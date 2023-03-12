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
					Slide     *slide
				},
			) error {
				params.Lifecycle.Append(
					fx.Hook{
						OnStart: func(ctx context.Context) error {
							err := params.Service.OnStart(ctx)
							if err != nil {
								return err
							}
							params.Service.SetConnectionListChange(params.Slide.SetConnectionListChange)
							params.Service.SetConnectionInstanceChange(params.Slide.SetConnectionInstanceChange)

							return nil
						},
						OnStop: func(ctx context.Context) error {
							return params.Service.OnStop(ctx)
						},
					},
				)
				return nil
			},
		),
	)
}
