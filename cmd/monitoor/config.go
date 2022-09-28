package main

import (
	"flag"
	"os"
	"time"
)

var cfg *Config = &Config{}

type Config struct {
	monitorTime time.Duration
	captureTime time.Duration
	logLevel    int
	persist     bool
	db          struct {
		driver       string
		dsn          string
		maxIdleConns int
		maxOpenConns int
		maxIdleTime  int
	}
}

func init() {
	const INFO_LEVEL int = 1

	flag.DurationVar(&cfg.monitorTime, "monitor-time", 1*time.Second, "monitor time")
	flag.DurationVar(&cfg.captureTime, "capture-time", 1*time.Hour, "capture time")

	flag.StringVar(&cfg.db.driver, "driver", os.Getenv("DB_DRIVER"), "database driver")
	flag.StringVar(&cfg.db.dsn, "dsn", os.Getenv("DB_DSN"), "database dsn")
	flag.IntVar(&cfg.db.maxIdleConns, "max-idle-conns", 5, "max idle connections")
	flag.IntVar(&cfg.db.maxOpenConns, "max-open-conns", 10, "max open connections")
	flag.IntVar(&cfg.db.maxIdleTime, "max-idle-time", 2, "max idle time")

	flag.IntVar(&cfg.logLevel, "log-level", INFO_LEVEL, "log level")
	flag.BoolVar(&cfg.persist, "persist", false, "enable persistance")

	flag.Parse()
}
