package model

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var ErrNoRows = errors.New("no rows for requested query")
var ErrTimedOut = errors.New("query time limit exceeded")

type Stat struct {
	Sent     uint64
	Received uint64
	Total    uint64
}

type Snapshot struct {
	Timestamp int64
	Stat
}

type MonthStat struct {
	Month string
	Stat
}

type DateStat struct {
	HoursMonitored int
	Date           string // YYYY-MM-DD
	Stat
}

type SnapshotModel struct {
	db *sql.DB
}

func NewSnapshotModel(db *sql.DB) *SnapshotModel {
	return &SnapshotModel{db: db}
}

func (m *SnapshotModel) Insert(ctx context.Context, s *Snapshot) error {
	query := `INSERT INTO snapshots (timestamp, sent, received, total) VALUES (?, ?, ?, ?)`

	timeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	args := []interface{}{s.Timestamp, s.Stat.Sent, s.Stat.Received, s.Stat.Total}

	_, err := m.db.ExecContext(timeout, query, args...)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return ErrTimedOut
		}

		return err
	}

	return nil
}
