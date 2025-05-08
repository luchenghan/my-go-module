package repository

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		fx.Annotate(
			nil,
			fx.ResultTags(`name:"master"`),
		),
	),
	fx.Provide(
		fx.Annotate(
			nil,
			fx.ResultTags(`name:"slave"`),
		),
	),
)
