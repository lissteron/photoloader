package elog

import (
	"context"

	"go.uber.org/zap"
)

type Logger interface {
	Fatalf(ctx context.Context, template string, args ...any)
	Fatal(ctx context.Context, args ...any)
	Errorw(ctx context.Context, msg string, args ...any)
	Errorf(ctx context.Context, template string, args ...any)
	Error(ctx context.Context, args ...any)
	Warnw(ctx context.Context, msg string, args ...any)
	Warnf(ctx context.Context, template string, args ...any)
	Warn(ctx context.Context, args ...any)
	Infof(ctx context.Context, template string, args ...any)
	Info(ctx context.Context, args ...any)
	Debugf(ctx context.Context, template string, args ...any)
	Debug(ctx context.Context, args ...any)

	With(args ...any) Logger
	SetLevel(lvl string)
	Sync()

	SetMeter(client Meter) error

	AtomicLevel() zap.AtomicLevel
	Zap() *zap.Logger
}

type Meter interface {
	IncByLevel(level string)
}
