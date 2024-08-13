package main

import (
	"context"
	"log"
	"os"
	"runtime/debug"

	"github.com/urfave/cli/v2"

	"github.com/lissteron/photoloader/cmd/photoloader/server"
	"github.com/lissteron/photoloader/config"
	"github.com/lissteron/photoloader/pkg/elog"
)

const (
	exitCodeOk = iota
	exitCodeAppError
	exitCodeLoggerError
	exitCodePanic
)

func main() {
	exitCode := exitCodeOk
	defer func() { os.Exit(exitCode) }()

	logger, err := elog.New(config.DefaultLogLevel)
	if err != nil {
		exitCode = exitCodeLoggerError

		panic(err)
	}

	defer logger.Sync()

	defer func() {
		if r := recover(); r != nil {
			logger.With("stack_trace", string(debug.Stack())).Errorf(context.Background(), "service panic: %v", r)

			exitCode = exitCodePanic
		}
	}()

	newApp := &cli.App{
		Commands: []*cli.Command{
			server.BuildCmd(),
			ver(),
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := newApp.RunContext(ctx, os.Args); err != nil {
		logger.Errorf(ctx, "run failed: %v", err)

		exitCode = exitCodeAppError
	}
}

func ver() *cli.Command {
	return &cli.Command{
		Name:        "version",
		Aliases:     []string{"ver", "v"},
		Description: "Show build info",
		Action: func(_ *cli.Context) error {
			log.Printf("ServiceName: %s", config.ServiceName)
			log.Printf("AppName: %s", config.AppName)
			log.Printf("GitHash: %s", config.GitHash)
			log.Printf("Version: %s", config.Version)
			log.Printf("BuildAt: %s", config.BuildAt)
			log.Printf("ReleaseID: %s", config.ReleaseID)

			return nil
		},
	}
}
