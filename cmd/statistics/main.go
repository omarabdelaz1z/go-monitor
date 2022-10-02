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
	"github.com/rs/zerolog"
)

func main() {
	var (
		db   *sql.DB
		err  error
		file *os.File
		cfg  *config.Config = &config.Config{}
	)

	flag.StringVar(&cfg.Db.Driver, "driver", os.Getenv("DB_DRIVER"), "database driver")
	flag.StringVar(&cfg.Db.Dsn, "dsn", os.Getenv("DB_DSN"), "database dsn")
	flag.IntVar(&cfg.Db.MaxIdleConns, "max-idle-conns", 5, "max idle connections")
	flag.IntVar(&cfg.Db.MaxOpenConns, "max-open-conns", 10, "max open connections")
	flag.IntVar(&cfg.Db.MaxIdleTime, "max-idle-time", 2, "max idle time")

	flag.StringVar(&cfg.Log.Path, "log-path", os.Getenv("LOG_PATH"), "log path")

	helper.EnumFlag(&cfg.Log.Level, "log-level", []string{"debug", "info", "warn", "error"}, "log level")
	flag.Parse()

	logLevel := zerolog.Level(helper.GetLevel(cfg.Log.Level))

	file, err = os.OpenFile(cfg.Log.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

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

	logger.Debug().Msg("initiating connection to database")

	db, err = provider.NewDatabase(&provider.DbConfig{
		Driver:       cfg.Db.Driver,
		Dsn:          cfg.Db.Dsn,
		MaxIdleConns: cfg.Db.MaxIdleConns,
		MaxOpenConns: cfg.Db.MaxOpenConns,
		MaxIdleTime:  cfg.Db.MaxIdleTime,
	})

	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initiate connection to database")
		return
	}

	logger.Info().Msg("connected to database")

	defer db.Close()

	service := &Service{
		snapshots:     model.NewSnapshotModel(db),
		config:        cfg,
		logger:        logger,
		monthSafeList: []string{},
	}

	if err = service.Run(); err != nil {
		logger.Fatal().Err(err).Msg("error occrred while running service")
		return
	}

	logger.Info().Msg("service stopped successfully")
}
