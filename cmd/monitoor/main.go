package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/omarabdelaz1z/go-monitor/cmd/provider"
	"github.com/omarabdelaz1z/go-monitor/internal/logger"
	"github.com/omarabdelaz1z/go-monitor/internal/model"
	m "github.com/omarabdelaz1z/go-monitor/pkg/monitoor"
)

var noprops = map[string]string{}

func main() {
	var (
		db           *sql.DB
		err          error
		snapshots    *model.SnapshotModel
		periodicStat *m.NetStat
	)

	logger := logger.New(os.Stdout, logger.LogLevel(cfg.logLevel))

	logger.Info("config loaded", map[string]string{
		"monitor_time": cfg.monitorTime.String(),
		"capture_time": cfg.captureTime.String(),
		"log_level":    fmt.Sprint(cfg.logLevel),
		"persistance":  fmt.Sprint(cfg.persist),
	})

	if cfg.persist {
		db, err = provider.NewDatabase(&provider.DbConfig{
			Driver:       cfg.db.driver,
			Dsn:          cfg.db.dsn,
			MaxIdleConns: cfg.db.maxIdleConns,
			MaxOpenConns: cfg.db.maxOpenConns,
			MaxIdleTime:  cfg.db.maxIdleTime,
		})

		if err != nil {
			logger.Fatal(err, map[string]string{
				"dsn":    cfg.db.dsn,
				"driver": cfg.db.driver,
				"action": "connecting to database",
			})
			return
		}

		defer db.Close()
		logger.Info("connected to database", noprops)

		snapshots = model.NewSnapshotModel(db)

		periodicStat = &m.NetStat{
			BytesSent:  0,
			BytesRecv:  0,
			BytesTotal: 0,
		}
	}

	service := &Service{
		config:    cfg,
		logger:    logger,
		snapshots: snapshots,

		monitorTicker: time.NewTicker(cfg.monitorTime),
		captureTicker: time.NewTicker(cfg.captureTime),

		cumulativeStat: &m.NetStat{
			BytesSent:  0,
			BytesRecv:  0,
			BytesTotal: 0,
		},

		periodicStat: periodicStat,
	}

	if service.Run(); err != nil {
		logger.Fatal(err, noprops)
	}

	logger.Info("service stopped successfully", noprops)
}
