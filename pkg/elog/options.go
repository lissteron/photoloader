package elog

import (
	"io"
	"sync"

	"github.com/lissteron/photoloader/pkg/helpers"
)

type options struct {
	callerSkip int
	meterMu    sync.RWMutex
	meter      Meter

	appVersion string
	appName    string

	output io.Writer
}

type Option interface {
	apply(o *options) error
}

type funcOption func(o *options) error

func (f funcOption) apply(o *options) error {
	return f(o)
}

func WithCallerSkip(skip int) Option {
	return funcOption(func(o *options) error {
		o.callerSkip = skip

		return nil
	})
}

func WithOutput(output io.Writer) Option {
	return funcOption(func(o *options) error {
		o.output = output

		return nil
	})
}

func WithAppName(appName string) Option {
	return funcOption(func(o *options) error {
		o.appName = appName

		return nil
	})
}

func WithAppVersion(appVersion string) Option {
	return funcOption(func(o *options) error {
		o.appVersion = appVersion

		return nil
	})
}

func WithMeter(meter Meter) Option {
	return funcOption(func(o *options) error {
		return o.setMeter(meter)
	})
}

func (o *options) setMeter(meter Meter) error {
	if helpers.IsNotNil(o.meter) {
		return nil
	}

	o.meter = meter

	return nil
}
