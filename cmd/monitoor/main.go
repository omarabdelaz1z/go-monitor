package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/omarabdelaz1z/go-monitor/cmd/config"
	"github.com/omarabdelaz1z/go-monitor/cmd/helper"
	"github.com/omarabdelaz1z/go-monitor/cmd/provider"
	"github.com/omarabdelaz1z/go-monitor/internal/model"
	m "github.com/omarabdelaz1z/go-monitor/pkg/monitoor"
	"github.com/rs/zerolog"
)

type monitoorConfig struct {
	monitorTime, captureTime time.Duration

	base *config.Config

	allowPersist bool
}

func main() {
	var (
		mCfg *monitoorConfig = &monitoorConfig{
			base:         &config.Config{},
			allowPersist: false,
			monitorTime:  time.Duration(0),
			captureTime:  time.Duration(0),
		}
	)

	flag.StringVar(&mCfg.base.Db.Driver, "driver", os.Getenv("DB_DRIVER"), "database driver")
	flag.StringVar(&mCfg.base.Db.Dsn, "dsn", os.Getenv("DB_DSN"), "database dsn")
	flag.IntVar(&mCfg.base.Db.MaxIdleConns, "max-idle-conns", 5, "max idle connections")
	flag.IntVar(&mCfg.base.Db.MaxOpenConns, "max-open-conns", 10, "max open connections")
	flag.IntVar(&mCfg.base.Db.MaxIdleTime, "max-idle-time", 2, "max idle time")
	flag.StringVar(&mCfg.base.Log.Path, "log-path", os.Getenv("LOG_PATH"), "log path")
	helper.EnumFlag(&mCfg.base.Log.Level, "log-level", []string{"debug", "info", "warn", "error"}, "log level")

	flag.DurationVar(&mCfg.captureTime, "capture-time", time.Hour*1, "Capture time")
	flag.DurationVar(&mCfg.monitorTime, "monitor-time", time.Second*1, "Monitor time")
	flag.BoolVar(&mCfg.allowPersist, "persist", false, "Persist data to database")

	flag.Parse()

	var (
		db           *sql.DB
		err          error
		snapshots    *model.SnapshotModel
		periodicStat *m.NetStat
	)

	logLevel := zerolog.Level(helper.GetLevel(mCfg.base.Log.Level))

	file, err := os.OpenFile(mCfg.base.Log.Path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

	if err != nil {
		panic(err)
	}

	logger := zerolog.New(zerolog.MultiLevelWriter(file, zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC1123,
		FormatCaller: func(i interface{}) string {
			if i == nil {
				return ""
			}
			return filepath.Base(fmt.Sprintf("%+v", i))
		},
	})).Level(logLevel).With().Timestamp().Logger()

	logger.Info().Msg("loggers initialized")
	logger.Info().Msg("config loaded")

	if mCfg.allowPersist {
		logger.Info().Msg("persistance allowed")
		logger.Debug().Caller().Msg("initiating connection to database")

		db, err = provider.NewDatabase(&provider.DbConfig{
			Driver:       mCfg.base.Db.Driver,
			Dsn:          mCfg.base.Db.Dsn,
			MaxIdleConns: mCfg.base.Db.MaxIdleConns,
			MaxOpenConns: mCfg.base.Db.MaxOpenConns,
			MaxIdleTime:  mCfg.base.Db.MaxIdleTime,
		})

		if err != nil {
			logger.Fatal().Err(err).Msg("failed to initiate connection to database")
			return
		}

		defer db.Close()
		logger.Info().Msg("connected to database")

		snapshots = model.NewSnapshotModel(db)

		periodicStat = &m.NetStat{
			BytesSent:  0,
			BytesRecv:  0,
			BytesTotal: 0,
		}
	}

	service := &Service{
		config:    mCfg,
		logger:    logger,
		snapshots: snapshots,

		monitorTicker: time.NewTicker(mCfg.monitorTime),
		captureTicker: time.NewTicker(mCfg.captureTime),

		cumulativeStat: &m.NetStat{
			BytesSent:  0,
			BytesRecv:  0,
			BytesTotal: 0,
		},

		periodicStat: periodicStat,
	}

	if err = service.Run(); err != nil {
		logger.Fatal().Err(err).Msg("error occrred while running service")
		return
	}

	logger.Info().Msg("service stopped successfully")
}
