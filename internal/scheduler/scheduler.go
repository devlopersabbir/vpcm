package scheduler

import (
	"log/slog"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron *cron.Cron
}

func New() *Scheduler {
	return &Scheduler{
		cron: cron.New(cron.WithSeconds()),
	}
}

func (s *Scheduler) AddJob(spec string, cmd func()) (int, error) {
	slog.Debug("Registering scheduled job", "spec", spec)
	id, err := s.cron.AddFunc(spec, cmd)
	return int(id), err
}

func (s *Scheduler) Start() {
	slog.Info("Starting cron scheduler")
	s.cron.Start()
}

func (s *Scheduler) Stop() {
	slog.Info("Stopping cron scheduler")
	s.cron.Stop()
}
