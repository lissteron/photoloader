package internal

import (
	"context"
	"time"
)

type Meter interface {
	CallDurationObserve(host, method, path string, status int, isErr bool, duration time.Duration)
	PanicInc(location string)
}

type Logger interface {
	Warn(ctx context.Context, args ...any)
	Error(ctx context.Context, args ...any)
}
