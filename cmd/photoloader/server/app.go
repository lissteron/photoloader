package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/lissteron/photoloader/config"
	"github.com/lissteron/photoloader/internal/app/core/services"
	"github.com/lissteron/photoloader/internal/app/handlers"
	"github.com/lissteron/photoloader/internal/app/repositories"
	"github.com/lissteron/photoloader/pkg/ehttp"
	"github.com/lissteron/photoloader/pkg/elog"
	"github.com/lissteron/photoloader/pkg/manager"
	"github.com/lissteron/photoloader/ui"
)

type App struct {
	logger elog.Logger
	config *config.Config

	manager manager.Manager
}

func NewApp(_ context.Context, logger elog.Logger, cfg *config.Config) (*App, error) {
	app := &App{
		logger:  logger,
		config:  cfg,
		manager: manager.New(),
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validate: %w", err)
	}

	httpServer := ehttp.NewServer(cfg.HTTPServer, ehttp.WithServerLogger(logger))

	// repositories
	photoRepo := repositories.NewPhoto(logger, cfg.PhotoConfig)

	// services
	photoService := services.NewPhoto(logger, photoRepo)

	// handlers
	photoHandler := handlers.NewPhoto(logger, photoService, cfg.PhotoConfig.Path)

	photoHandler.SetRounte(httpServer.Router())
	httpServer.Router().Handle("/", ui.NewFS())

	app.manager.Add("http server", httpServer)

	return app, nil
}

func (a *App) Run(ctx context.Context) error {
	a.logger.Info(ctx, "start application")

	// start application.
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	managerCtx := a.manager.Run(ctx)

	// wait exit signal
	select {
	case <-exit:
		a.logger.Info(ctx, "stopping application")
	case <-managerCtx.Done():
		a.logger.Error(ctx, "stopping application with error")
	case <-ctx.Done():
		a.logger.Error(ctx, "stopping application with context canceled")
	}

	signal.Stop(exit)

	return nil
}

func (a *App) Close(ctx context.Context) {
	if err := a.manager.Stop(ctx); err != nil {
		a.logger.Errorf(ctx, "manager stop: %v", err)
	}

	a.logger.Info(ctx, "service exited")
}
