package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Stat struct {
	Sent     uint64
	Received uint64
	Total    uint64
}

type Snapshot struct {
	Timestamp int64
	Stat      Stat
}

type MonthStat struct {
	Month string
	Stat  Stat
}

type SnapshotModel struct {
	db *sql.DB
}

func NewSnapshotModel(db *sql.DB) *SnapshotModel {
	return &SnapshotModel{db: db}
}

func (m *SnapshotModel) Insert(s *Snapshot) error {
	query := `INSERT INTO snapshots (timestamp, sent, received, total) VALUES (?, ?, ?, ?)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()
	args := []interface{}{s.Timestamp, s.Stat.Sent, s.Stat.Received, s.Stat.Total}

	_, err := m.db.ExecContext(ctx, query, args...)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("inserting snapshot timed out %v", err)
		}

		return fmt.Errorf("failed to insert snapshot: %s", err)
	}

	return nil
}

func (m *SnapshotModel) GetStatsByMonth(month string) ([]Snapshot, error) {
	query := `SELECT 
		strftime('%s', strftime('%Y-%m-%d', timestamp, 'unixepoch')) AS day_unix,
		SUM(sent),
		SUM(received),
		SUM(total)
	FROM snapshots
	WHERE strftime('%m', timestamp, 'unixepoch') = ?
	GROUP BY day_unix
	ORDER BY day_unix DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	rows, err := m.db.QueryContext(ctx, query, month)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("querying snapshots timed out %v", err)
		}

		return nil, fmt.Errorf("failed to query: %s", err)
	}

	var stats []Snapshot

	defer rows.Close()

	for rows.Next() {
		var s Snapshot

		err := rows.Scan(&s.Timestamp, &s.Stat.Sent, &s.Stat.Received, &s.Stat.Total)

		if err != nil {
			return stats, fmt.Errorf("found no rows: %s", err)
		}

		stats = append(stats, s)
	}

	return stats, nil
}

func (m *SnapshotModel) GetMonthStat(month string) (MonthStat, error) {
	query := `SELECT 
		strftime('%m', timestamp, 'unixepoch') AS month, SUM(sent), SUM(received), SUM(total) 
		FROM snapshots 
		WHERE month = ?
		GROUP BY month`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var s MonthStat

	if err := m.db.QueryRowContext(ctx, query, month).Scan(&s.Month, &s.Stat.Sent, &s.Stat.Received, &s.Stat.Total); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return s, fmt.Errorf("querying snapshot timed out %v", err)
		}
		if errors.Is(err, sql.ErrNoRows) {
			return s, fmt.Errorf("found no rows: %s", err)
		}

		return s, fmt.Errorf("failed to query: %s", err)
	}

	return s, nil
}

func (m *SnapshotModel) GetAllStats() ([]Snapshot, error) {
	query := `SELECT 
					strftime('%s', strftime('%Y-%m-%d', timestamp, 'unixepoch')) AS day_unix, SUM(sent), SUM(received), SUM(total)
					FROM snapshots
					GROUP BY day_unix
					ORDER BY day_unix DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	rows, err := m.db.QueryContext(ctx, query)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("querying snapshots timed out %v", err)
		}

		return nil, fmt.Errorf("failed to query: %s", err)
	}

	var stats []Snapshot

	defer rows.Close()

	for rows.Next() {
		var s Snapshot

		err := rows.Scan(&s.Timestamp, &s.Stat.Sent, &s.Stat.Received, &s.Stat.Total)

		if err != nil {
			return stats, fmt.Errorf("found no rows: %s", err)
		}

		stats = append(stats, s)
	}

	return stats, nil
}
