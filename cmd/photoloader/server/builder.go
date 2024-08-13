package server

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/lissteron/photoloader/config"
	"github.com/lissteron/photoloader/pkg/elog"
)

func BuildCmd() *cli.Command {
	cfg := config.NewConfig()

	return &cli.Command{
		Name:        "server",
		Description: "start http service",
		Action: func(ctx *cli.Context) error {
			return run(ctx.Context, cfg)
		},
	}
}

func run(ctx context.Context, cfg *config.Config) error {
	logger, err := elog.New(cfg.Base.LogLevel)
	if err != nil {
		return fmt.Errorf("init elog: %w", err)
	}

	defer logger.Sync()

	app, err := NewApp(ctx, logger, cfg)
	if err != nil {
		return fmt.Errorf("new app: %w", err)
	}

	defer app.Close(ctx)

	if err := app.Run(ctx); err != nil {
		return fmt.Errorf("app run: %w", err)
	}

	return nil
}
