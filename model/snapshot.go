package model

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db  *sql.DB
	err error
)

func InitDb(driver, dsn string) error {
	db, err = sql.Open(driver, dsn)

	if err != nil {
		return fmt.Errorf("failed to open database: %s", err)
	}

	return db.Ping()
}

type Snapshot struct {
	Timestamp int64  `json:"timestamp"`
	Sent      uint64 `json:"sent"`
	Received  uint64 `json:"received"`
	Total     uint64 `json:"total"`
}

func Insert(stat *Snapshot) error {
	_, err = db.Exec(
		"INSERT INTO snapshots (timestamp, sent, received, total) VALUES (?, ?, ?, ?)",
		stat.Timestamp,
		stat.Sent,
		stat.Received,
		stat.Total,
	)

	if err != nil {
		return fmt.Errorf("failed to insert snapshot: %s", err)
	}

	return nil
}
