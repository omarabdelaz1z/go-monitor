package main

import (
	"fmt"
	"os"
	"time"

	"github.com/omarabdelaz1z/go-monitor/internal/logger"
	"github.com/omarabdelaz1z/go-monitor/internal/model"
	m "github.com/omarabdelaz1z/go-monitor/pkg/monitoor"
)

var noprops = map[string]string{}

func main() {
	logger := logger.New(os.Stdout, logger.LogLevel(cfg.logLevel))

	logger.Info("config loaded", map[string]string{
		"monitor_time": cfg.monitorTime.String(),
		"capture_time": cfg.captureTime.String(),
		"log_level":    fmt.Sprint(cfg.logLevel),
	})

	db, err := NewDatabase(cfg)

	if err != nil {
		logger.Error(err, map[string]string{
			"dsn":    cfg.db.dsn,
			"driver": cfg.db.driver,
			"action": "connecting to database",
		})
	}

	logger.Info("connected to database", noprops)

	defer db.Close()

	service := &Service{
		config:    cfg,
		logger:    logger,
		snapshots: model.NewSnapshotModel(db),

		monitorTicker: time.NewTicker(cfg.monitorTime),
		captureTicker: time.NewTicker(cfg.captureTime),

		cumulativeStat: &m.NetStat{
			BytesSent:  0,
			BytesRecv:  0,
			BytesTotal: 0,
		},

		periodicStat: &m.NetStat{
			BytesSent:  0,
			BytesRecv:  0,
			BytesTotal: 0,
		},
	}

	if service.Run(); err != nil {
		logger.Fatal(err, noprops)
	}

	logger.Info("service stopped successfully", noprops)
}
