package main

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/minilik/ecommerce/config"
	di_container "github.com/minilik/ecommerce/internal/infrastructure/container"
)

func main() {
	cfg, err := config.Load("config")
	if err != nil {
		panic(fmt.Errorf("load config: %w", err))
	}

	app, err := di_container.Build(cfg)
	if err != nil {
		panic(fmt.Errorf("build container: %w", err))
	}
	defer func() {
		if err := app.Close(); err != nil {
			app.Logger.Error("failed to close container", zap.Error(err))
		}
	}()

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	app.Logger.Info("starting HTTP server", zap.String("address", addr))

	if err := app.Router.Run(addr); err != nil {
		app.Logger.Fatal("server exited with error", zap.Error(err))
	}
}
