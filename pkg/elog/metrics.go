package elog

import "github.com/lissteron/photoloader/pkg/helpers"

const (
	_levelDebug = "debug"
	_levelInfo  = "info"
	_levelWarn  = "warn"
	_levelError = "error"
	_levelFatal = "fatal"
)

func (l *LoggerPKG) incr(level string) {
	if lockOk := l.meterMu.TryRLock(); !lockOk {
		return
	}

	defer l.meterMu.RUnlock()

	if helpers.IsNotNil(l.meter) {
		l.meter.IncByLevel(level)
	}
}
