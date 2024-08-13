package manager

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

type job struct {
	name string
	job  Job
}

//nolint:revive // lint mistake
type ManagerPKG struct {
	errGroup    *errgroup.Group
	errGroupCtx context.Context //nolint:containedctx // need context in struct.
	jobs        []job
}

func New() *ManagerPKG {
	return &ManagerPKG{}
}

func (m *ManagerPKG) Add(name string, j Job) {
	m.jobs = append(m.jobs, job{
		name: name,
		job:  j,
	})
}

func (m *ManagerPKG) Run(ctx context.Context) context.Context {
	// error group
	m.errGroup, m.errGroupCtx = errgroup.WithContext(ctx)

	for _, j := range m.jobs {
		m.errGroup.Go(func() error {
			if err := j.job.Start(ctx); err != nil {
				return fmt.Errorf("job %s start: %w", j.name, err)
			}

			return nil
		})
	}

	return m.errGroupCtx
}

func (m *ManagerPKG) Stop(ctx context.Context) error {
	const (
		closeTimeout = 2 * time.Second
		exitTimeout  = 30 * time.Second
	)

	// wait for orchestration
	time.Sleep(closeTimeout)

	ctx, cancel := context.WithTimeout(ctx, exitTimeout)
	defer cancel()

	for i := len(m.jobs); i > 0; i-- {
		if err := m.jobs[i-1].job.Stop(ctx); err != nil {
			return fmt.Errorf("job %s stop: %w", m.jobs[i-1].name, err)
		}
	}

	if err := m.errGroup.Wait(); err != nil {
		return fmt.Errorf("err group wait: %w", err)
	}

	return nil
}
