package main

import (
	"context"
	"fmt"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/omarabdelaz1z/go-monitor/cmd/monitoor/helper"
	"github.com/omarabdelaz1z/go-monitor/internal/logger"
	"github.com/omarabdelaz1z/go-monitor/internal/model"
	"github.com/omarabdelaz1z/go-monitor/internal/util"
	m "github.com/omarabdelaz1z/go-monitor/pkg/monitoor"
	"golang.org/x/sync/errgroup"
)

type Service struct {
	logger *logger.Logger
	config *Config

	snapshots *model.SnapshotModel

	monitorTicker, captureTicker *time.Ticker
	cumulativeStat, periodicStat *m.NetStat

	mu sync.RWMutex
}

func (s *Service) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	g, gCtx := errgroup.WithContext(ctx)

	buffer := make(chan *m.NetStat)

	g.Go(func() error {
		s.logger.Info("monitor started", nil)
		return s.Monitor(gCtx, buffer)
	})

	g.Go(func() error {
		s.logger.Info("display started", nil)
		return s.Display(gCtx, buffer)
	})

	g.Go(func() error {
		s.logger.Info("capture started", nil)
		return s.Capture(gCtx, buffer)
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("service failed while running: %w", err)
	}

	return nil
}

func (s *Service) Monitor(ctx context.Context, buffer chan<- *m.NetStat) error {
	var currentStat *m.NetStat
	var err error

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("monitor stopped", noprops)
			s.monitorTicker.Stop()
			close(buffer)
			return nil
		case <-s.monitorTicker.C:
			if nil == currentStat {
				currentStat, err = m.Brief()

				if err != nil {
					s.logger.Warn(fmt.Sprintf("failed to get current netstat %v", err), noprops)
					continue // retry again
				}
			}

			newStat, err := m.Brief()

			if err != nil {
				s.logger.Warn(fmt.Sprintf("failed to get new netstat %v", err), noprops)
				continue // retry again
			}

			delta := helper.Delta(newStat, currentStat)
			buffer <- delta

			s.mu.Lock()

			s.logger.Debug("update periodicStat stat", nil)
			helper.UpdateWith(s.periodicStat, helper.Incr(s.periodicStat, delta))

			s.logger.Debug("update cumulative stat", nil)
			helper.UpdateWith(s.cumulativeStat, helper.Incr(s.cumulativeStat, delta))

			s.mu.Unlock()

			s.logger.Debug("update cumulative stat", nil)
			helper.UpdateWith(currentStat, *newStat)
		}
	}
}

func (s *Service) Display(ctx context.Context, buffer <-chan *m.NetStat) error {
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("display stopped", nil)
			s.captureTicker.Stop()
			return nil
		case stat, ok := <-buffer:
			if !ok {
				s.logger.Info("display stopped", nil)
				return fmt.Errorf("buffer channel is closed")
			}

			s.mu.RLock()
			cumulative := util.ByteCountSI(s.cumulativeStat.BytesTotal)
			s.mu.RUnlock()

			s.logger.Info(
				"monitored",
				map[string]string{
					"sent":       util.ByteCountSI(stat.BytesSent),
					"received":   util.ByteCountSI(stat.BytesRecv),
					"total":      util.ByteCountSI(stat.BytesTotal),
					"cumulative": cumulative,
				},
			)
		}
	}
}

func (s *Service) Capture(ctx context.Context, buffer <-chan *m.NetStat) error {
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("capture stopped", nil)
			s.captureTicker.Stop()
			return nil
		case <-s.captureTicker.C:
			s.mu.RLock()

			snap := &model.Snapshot{
				Timestamp: time.Now().Unix(),
				Stat: model.Stat{
					Sent:     s.periodicStat.BytesSent,
					Received: s.periodicStat.BytesRecv,
					Total:    s.periodicStat.BytesTotal,
				},
			}

			s.logger.Debug("inserting stat into database", map[string]string{
				"sent":      fmt.Sprint(snap.Stat.Sent),
				"recv":      fmt.Sprint(snap.Stat.Received),
				"total":     fmt.Sprint(snap.Stat.Total),
				"timestamp": fmt.Sprint(snap.Timestamp),
			})

			s.snapshots.Insert(snap)
			s.mu.RUnlock()

			s.mu.Lock()
			s.logger.Debug("reset periodic stat", nil)
			helper.UpdateWith(s.periodicStat, m.NetStat{})
			s.mu.Unlock()
		}
	}
}
