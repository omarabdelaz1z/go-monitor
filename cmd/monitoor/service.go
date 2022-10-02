package main

import (
	"context"
	"fmt"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/omarabdelaz1z/go-monitor/cmd/monitoor/helper"
	"github.com/omarabdelaz1z/go-monitor/internal/model"
	"github.com/omarabdelaz1z/go-monitor/internal/util"
	m "github.com/omarabdelaz1z/go-monitor/pkg/monitoor"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

type Service struct {
	config    *monitoorConfig
	snapshots *model.SnapshotModel
	mu        sync.RWMutex

	logger zerolog.Logger

	monitorTicker, captureTicker *time.Ticker
	cumulativeStat, periodicStat *m.NetStat
}

func (s *Service) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	g, gCtx := errgroup.WithContext(ctx)

	buffer := make(chan *m.NetStat)

	g.Go(func() error {
		s.logger.Info().Msg("monitor goroutine launched")
		return s.Monitor(gCtx, buffer)
	})

	g.Go(func() error {
		s.logger.Info().Msg("display goroutine launched")
		return s.Display(gCtx, buffer)
	})

	if s.config.allowPersist {
		g.Go(func() error {
			s.logger.Info().Msg("capture goroutine launched")
			return s.Capture(gCtx, buffer)
		})
	}

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
			s.logger.Info().Msg("monitor stopped")
			s.monitorTicker.Stop()
			close(buffer)
			return nil
		case <-s.monitorTicker.C:
			if nil == currentStat {
				currentStat, err = m.Brief()

				if err != nil {
					s.logger.Warn().Caller().Err(err).Msg("failed to get current stat")
					continue // retry again
				}
			}
			var newStat *m.NetStat

			newStat, err = m.Brief()

			if err != nil {
				s.logger.Warn().Err(err).Msg("failed to get new netstat")
				continue // retry again
			}

			delta := helper.Delta(newStat, currentStat)
			buffer <- delta

			s.mu.Lock()
			if s.config.allowPersist {
				helper.UpdateWith(s.periodicStat, helper.Incr(s.periodicStat, delta))
			}

			helper.UpdateWith(s.cumulativeStat, helper.Incr(s.cumulativeStat, delta))
			s.mu.Unlock()

			helper.UpdateWith(currentStat, *newStat)
		}
	}
}

func (s *Service) Display(ctx context.Context, buffer <-chan *m.NetStat) error {
	for {
		select {
		case <-ctx.Done():
			s.logger.Info().Msg("display stopped")
			s.captureTicker.Stop()
			return nil
		case stat, ok := <-buffer:
			if !ok {
				s.logger.Error().Caller().Msg("buffer channel is closed")
				return fmt.Errorf("buffer channel is closed")
			}

			s.mu.RLock()
			cumulative := util.ByteCountSI(s.cumulativeStat.BytesTotal)
			s.mu.RUnlock()

			s.logger.Info().
				Str("service", "display").
				Str("sent", util.ByteCountSI(stat.BytesSent)).
				Str("received", util.ByteCountSI(stat.BytesRecv)).
				Str("total", util.ByteCountSI(stat.BytesTotal)).
				Str("cumulative", cumulative).
				Send()
		}
	}
}

func (s *Service) Capture(ctx context.Context, buffer <-chan *m.NetStat) error {
	for {
		select {
		case <-ctx.Done():
			s.logger.Info().Msg("capture stopped")
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

			s.logger.Debug().Fields(map[string]string{
				"sent":      fmt.Sprint(snap.Stat.Sent),
				"recv":      fmt.Sprint(snap.Stat.Received),
				"total":     fmt.Sprint(snap.Stat.Total),
				"timestamp": fmt.Sprint(snap.Timestamp),
			}).Msg("persisting snapshot")

			err := s.snapshots.Insert(ctx, snap)

			if err != nil {
				s.logger.Error().Caller().Err(err)

				s.mu.RUnlock()
				continue
			}

			s.mu.RUnlock()

			s.mu.Lock()
			helper.UpdateWith(s.periodicStat, m.NetStat{})
			s.mu.Unlock()
		}
	}
}
