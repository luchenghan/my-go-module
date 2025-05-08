package bundler

import (
	"mymodule/fx/repository"
	"mymodule/fx/service"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.NopLogger,
	repository.Module,
	service.Module,
)
