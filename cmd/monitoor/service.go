package main

import (
	"context"
	"sync"
	"time"

	"github.com/omarabdelaz1z/go-monitor/internal/model"
	m "github.com/omarabdelaz1z/go-monitor/pkg/monitoor"
)

type Service struct {
	// logger *logger.Logger
	config *Config

	snapshots *model.SnapshotModel

	monitorTicker, captureTicker *time.Ticker
	cumulativeStat, periodicStat *m.NetStat

	mu sync.RWMutex
}

func (s *Service) Run() error {
	return nil
}

func (s *Service) Monitor(ctx context.Context, buffer chan<- *m.NetStat) error {
	return nil
}

func (s *Service) Display(ctx context.Context, buffer <-chan *m.NetStat) error {
	return nil
}

func (s *Service) Capture(ctx context.Context, buffer <-chan *m.NetStat) error {
	return nil
}
