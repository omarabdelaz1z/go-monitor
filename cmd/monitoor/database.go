package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func NewDatabase(cfg *Config) (*sql.DB, error) {
	db, err := sql.Open(cfg.db.driver, cfg.db.dsn)

	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %s", err)
	}

	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetConnMaxIdleTime(time.Duration(cfg.db.maxIdleTime))

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %s", err)
	}

	return db, nil
}
