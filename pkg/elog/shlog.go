package elog

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"runtime/debug"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var ErrInvalidLevel = errors.New("invalid log level")

type LoggerPKG struct {
	level         zap.AtomicLevel
	sugaredLogger *zap.SugaredLogger

	*options
}

func New(lvl string, opts ...Option) (*LoggerPKG, error) {
	logger := &LoggerPKG{
		level:   zap.NewAtomicLevelAt(zapLevel(lvl)),
		options: &options{},
	}

	logger.findTags()

	for _, opt := range opts {
		if err := opt.apply(logger.options); err != nil {
			return nil, fmt.Errorf("apply option: %w", err)
		}
	}

	if logger.output == nil {
		logger.output = os.Stderr
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.Lock(zapcore.AddSync(logger.output)),
		logger.level,
	)

	zopts := []zap.Option{
		zap.AddCallerSkip(logger.callerSkip + 1),
		zap.AddCaller(),
	}

	logger.sugaredLogger = zap.New(core, zopts...).Sugar()

	if logger.appName != "" {
		logger.sugaredLogger = logger.sugaredLogger.With("app_name", logger.appName)
	}

	if logger.appVersion != "" {
		logger.sugaredLogger = logger.sugaredLogger.With("app_version", logger.appVersion)
	}

	return logger, nil
}

func (l *LoggerPKG) SetLevel(lvl string) {
	l.level.SetLevel(zapLevel(lvl))
}

func (l *LoggerPKG) AtomicLevel() zap.AtomicLevel {
	return l.level
}

func (l *LoggerPKG) Sync() {
	// Ignore error because of https://github.com/uber-go/zap/issues/328.
	_ = l.sugaredLogger.Sync()
}

func (l *LoggerPKG) Zap() *zap.Logger {
	return l.sugaredLogger.Desugar()
}

func zapLevel(lvl string) zapcore.Level {
	level, _ := parseLevel(lvl)

	return level
}

func parseLevel(lvl string) (zapcore.Level, error) {
	switch strings.ToLower(lvl) {
	case "debug":
		return zap.DebugLevel, nil
	case "info":
		return zap.InfoLevel, nil
	case "warn", "warning":
		return zap.WarnLevel, nil
	case "err", "error":
		return zap.ErrorLevel, nil
	default:
		return zap.InfoLevel, ErrInvalidLevel
	}
}

func (l *LoggerPKG) findTags() {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	for _, opt := range buildInfo.Settings {
		if opt.Key == "-ldflags" {
			kv := parseInfo(opt.Value)

			l.appName = kv["AppName"]
			l.appVersion = kv["Version"]

			return
		}
	}
}

func parseInfo(input string) map[string]string {
	var (
		reg  = regexp.MustCompile(`-X\s+\S+\.([a-zA-Z0-9]+)=(\S+)`)
		resp = make(map[string]string)
	)

	for _, v := range reg.FindAllStringSubmatch(input, -1) {
		resp[v[1]] = v[2]
	}

	return resp
}
