package config

import "time"

//nolint:gochecknoglobals // build params
var (
	ServiceName string
	AppName     string
	Version     string
	GitHash     string
	BuildAt     string
	ReleaseID   string
	StartTime   = time.Now()
)
