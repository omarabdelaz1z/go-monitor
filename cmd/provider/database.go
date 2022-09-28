package provider

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DbConfig struct {
	MaxIdleConns, MaxOpenConns, MaxIdleTime int
	Driver, Dsn                             string
}

func NewDatabase(cfg *DbConfig) (*sql.DB, error) {
	db, err := sql.Open(cfg.Driver, cfg.Dsn)

	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %s", err)
	}

	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetConnMaxIdleTime(time.Duration(cfg.MaxIdleTime))

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %s", err)
	}

	return db, nil
}
