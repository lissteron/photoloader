package manager

import "context"

type Job interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type Manager interface {
	Add(name string, j Job)
	Run(ctx context.Context) context.Context
	Stop(ctx context.Context) error
}
