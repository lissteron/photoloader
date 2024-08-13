package elog

import (
	"context"
	"strings"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const (
	_traceIDKey = "trace_id"
	_spanIDKey  = "span_id"
)

func (l *LoggerPKG) Fatalf(ctx context.Context, template string, args ...any) {
	l.incr(_levelFatal)

	l.withTrace(ctx).Fatalf(prepareTemplate(template), args...)
}

func (l *LoggerPKG) Fatal(ctx context.Context, args ...any) {
	l.incr(_levelFatal)

	l.withTrace(ctx).Fatal(args...)
}

func (l *LoggerPKG) Errorw(ctx context.Context, msg string, args ...any) {
	l.incr(_levelError)
	setErrorTag(ctx)

	l.withTrace(ctx).Errorw(msg, args...)
}

func (l *LoggerPKG) Errorf(ctx context.Context, template string, args ...any) {
	l.incr(_levelError)
	setErrorTag(ctx)

	l.withTrace(ctx).Errorf(prepareTemplate(template), args...)
}

func (l *LoggerPKG) Error(ctx context.Context, args ...any) {
	l.incr(_levelError)
	setErrorTag(ctx)

	l.withTrace(ctx).Error(args...)
}

func (l *LoggerPKG) Warnw(ctx context.Context, msg string, args ...any) {
	l.incr(_levelWarn)

	l.withTrace(ctx).Warnw(msg, args...)
}

func (l *LoggerPKG) Warnf(ctx context.Context, template string, args ...any) {
	l.incr(_levelWarn)

	l.withTrace(ctx).Warnf(prepareTemplate(template), args...)
}

func (l *LoggerPKG) Warn(ctx context.Context, args ...any) {
	l.incr(_levelWarn)

	l.withTrace(ctx).Warn(args...)
}

func (l *LoggerPKG) Infof(ctx context.Context, template string, args ...any) {
	l.incr(_levelInfo)

	l.withTrace(ctx).Infof(prepareTemplate(template), args...)
}

func (l *LoggerPKG) Info(ctx context.Context, args ...any) {
	l.incr(_levelInfo)

	l.withTrace(ctx).Info(args...)
}

func (l *LoggerPKG) Debugf(ctx context.Context, template string, args ...any) {
	l.incr(_levelDebug)

	l.withTrace(ctx).Debugf(prepareTemplate(template), args...)
}

func (l *LoggerPKG) Debug(ctx context.Context, args ...any) {
	l.incr(_levelDebug)

	l.withTrace(ctx).Debug(args...)
}

func (l LoggerPKG) With(args ...any) Logger {
	l.sugaredLogger = l.sugaredLogger.With(args...)

	return &l
}

func (l *LoggerPKG) SetMeter(meter Meter) error {
	l.meterMu.Lock()
	defer l.meterMu.Unlock()

	return l.options.setMeter(meter)
}

func prepareTemplate(template string) string {
	// %w verb can be used only in 'fmt.Errorf' calls.
	return strings.ReplaceAll(template, "%w", "%v")
}

func (l *LoggerPKG) withTrace(ctx context.Context) *zap.SugaredLogger {
	if sc := trace.SpanContextFromContext(ctx); sc.IsValid() {
		return l.sugaredLogger.Desugar().With(
			zap.Stringer(_traceIDKey, sc.TraceID()),
			zap.Stringer(_spanIDKey, sc.SpanID()),
		).Sugar()
	}

	return l.sugaredLogger
}

func setErrorTag(ctx context.Context) {
	if span := trace.SpanFromContext(ctx); span != nil {
		span.SetStatus(codes.Error, "error")
	}
}
