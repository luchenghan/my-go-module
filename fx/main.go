package main

import (
	"mymodule/fx/bundler"

	"go.uber.org/fx"
)

func main() {
	fx.New(bundler.Module).Run()
}
