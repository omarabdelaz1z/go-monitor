package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/omarabdelaz1z/go-monitor/cmd/config"
	"github.com/omarabdelaz1z/go-monitor/internal/model"
	"github.com/omarabdelaz1z/go-monitor/internal/util"

	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

type Service struct {
	snapshots *model.SnapshotModel
	config    *config.Config
	logger    zerolog.Logger

	monthSafeList []string
}

func (s *Service) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	g, gCtx := errgroup.WithContext(ctx)

	months, err := s.snapshots.GetMonthsInYear(gCtx, fmt.Sprint(time.Now().Year()))

	if err != nil {
		return fmt.Errorf("failed to get months: %w", err)
	}

	s.monthSafeList = transformMonthName(months...)

	g.Go(func() error {
		return s.Respond(gCtx)
	})

	if err = g.Wait(); err != nil {
		s.logger.Error().Err(err).Caller().Msg("service failed while running")
		return fmt.Errorf("service failed while running: %w", err)
	}

	return nil
}

func (s *Service) Respond(ctx context.Context) error {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			option, _, err := selectPrompt("What would you like to do?", "View today's stats", "View stats for a month", "View all stats", "Exit")

			if option == -1 && nil == err {
				return nil
			}

			if err != nil {
				return err
			}

			switch option {
			case 0:
				err = s.HandleTodayStats(ctx, t)

				if err != nil {
					switch err {
					case model.ErrNoRows:
						fmt.Println("No stats for today")
					case model.ErrTimedOut:
						fmt.Println("Timed out while fetching stats")
					default:
						return err
					}
				}

				t.Render()
			case 1:
				err = s.HandleMonthStats(ctx, t)

				if err != nil {
					return err
				}

				t.Render()
			case 2:
				err = s.HandleAllStats(ctx, t)

				if err != nil {
					return err
				}

				t.Render()
			case 3:
				return nil
			}

			t.ResetHeaders()
			t.ResetRows()
			t.ResetFooters()
		}
	}
}

func (s *Service) HandleMonthStats(ctx context.Context, t table.Writer) error {
	var (
		err        error
		option     string
		dailyStats []model.Snapshot
	)

	_, option, err = selectPrompt("Which month would you like to view?", s.monthSafeList...)

	if err != nil {
		return err
	}

	if nil == err && option == "" {
		return nil
	}

	dailyStats, err = s.snapshots.GetStatsByMonth(ctx, option[:2])

	if err != nil {
		return fmt.Errorf("failed to get stats by month: %w", err)
	}

	t.SetCaption(fmt.Sprintf("Stats for %s", option[5:]))
	t.AppendHeader(table.Row{"Date", "Uploaded", "Downloaded", "Total"})

	for i := 0; i < len(dailyStats)-1; i++ {
		t.AppendRow(table.Row{
			time.Unix(dailyStats[i].Timestamp, 0).Format("2006-01-02"),
			util.ByteCountSI(dailyStats[i].Stat.Sent),
			util.ByteCountSI(dailyStats[i].Stat.Received),
			util.ByteCountSI(dailyStats[i].Stat.Total),
		})
	}

	t.AppendSeparator()
	t.AppendFooter(table.Row{
		"Cumulative",
		util.ByteCountSI(dailyStats[len(dailyStats)-1].Stat.Sent),
		util.ByteCountSI(dailyStats[len(dailyStats)-1].Stat.Received),
		util.ByteCountSI(dailyStats[len(dailyStats)-1].Stat.Total),
	})

	return nil
}

func (s *Service) HandleAllStats(ctx context.Context, t table.Writer) error {
	stats, err := s.snapshots.GetAllStats(ctx)

	if err != nil {
		return fmt.Errorf("failed to get stats: %w", err)
	}

	t.SetCaption("All stats")
	t.AppendHeader(table.Row{"Month", "Uploaded", "Downloaded", "Total"})

	for _, stat := range stats {
		t.AppendRow(table.Row{
			time.Unix(stat.Timestamp, 0).Format("2006-01-02"),
			util.ByteCountSI(stat.Stat.Sent),
			util.ByteCountSI(stat.Stat.Received),
			util.ByteCountSI(stat.Stat.Total),
		})
	}

	return nil
}

func (s *Service) HandleTodayStats(ctx context.Context, t table.Writer) error {
	today := time.Now().Format("2006-01-02")

	stat, err := s.snapshots.GetStatByDate(ctx, today)

	if err != nil {
		return err
	}

	t.SetCaption(fmt.Sprintf("Monitored %d hours on %s", stat.HoursMonitored, today))
	t.AppendHeader(table.Row{"Uploaded", "Downloaded", "Total"})

	t.AppendRow(table.Row{
		util.ByteCountSI(stat.Stat.Sent),
		util.ByteCountSI(stat.Stat.Received),
		util.ByteCountSI(stat.Stat.Total),
	})

	return nil
}
