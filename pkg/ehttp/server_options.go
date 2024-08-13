package ehttp

import (
	"time"
)

type ServerOption interface {
	apply(s *ServerPKG)
}

type funcServerOption func(s *ServerPKG)

func (f funcServerOption) apply(s *ServerPKG) {
	f(s)
}

func WithServerMeter(meter Meter) ServerOption {
	return funcServerOption(func(s *ServerPKG) {
		s.meter = meter
	})
}

func WithServerLogger(logger Logger) ServerOption {
	return funcServerOption(func(s *ServerPKG) {
		s.logger = logger
	})
}

func WithServerSlowRequestLog(
	logger Logger,
	duration time.Duration,
) ServerOption {
	return funcServerOption(func(s *ServerPKG) {
		if duration > 0 {
			s.logger = logger
			s.slowRequestDuration = duration
		}
	})
}
