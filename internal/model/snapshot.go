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

func (m *SnapshotModel) GetStatsByMonth(ctx context.Context, month string) ([]Snapshot, error) {
	query := `
	SELECT strftime('%s', strftime('%Y-%m-%d', timestamp, 'unixepoch', 'localtime')) AS unix,
		SUM(sent),
		SUM(received),
		SUM(total)
	FROM snapshots
	WHERE strftime('%m', timestamp, 'unixepoch', 'localtime') = ?
	GROUP BY unix

	UNION

	SELECT strftime('%m', timestamp, 'unixepoch', 'localtime') AS unix, SUM(sent), SUM(received), SUM(total) 
	FROM snapshots 
	WHERE unix = ?
	GROUP BY unix

	ORDER BY unix DESC`

	timeout, cancel := context.WithTimeout(ctx, 3*time.Second)

	defer cancel()

	rows, err := m.db.QueryContext(timeout, query, month, month)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, ErrTimedOut
		}

		return nil, err
	}

	var stats []Snapshot

	defer rows.Close()

	for rows.Next() {
		var s Snapshot

		if err = rows.Scan(&s.Timestamp, &s.Stat.Sent, &s.Stat.Received, &s.Stat.Total); err != nil {
			return nil, err
		}

		stats = append(stats, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}

func (m *SnapshotModel) GetMonthStat(ctx context.Context, month string) (MonthStat, error) {
	query := `SELECT 
		strftime('%m', timestamp, 'unixepoch') AS month, SUM(sent), SUM(received), SUM(total) 
		FROM snapshots 
		WHERE month = ?
		GROUP BY month`

	timeout, cancel := context.WithTimeout(ctx, 3*time.Second)

	defer cancel()

	var s MonthStat

	if err := m.db.QueryRowContext(timeout, query, month).Scan(&s.Month, &s.Stat.Sent, &s.Stat.Received, &s.Stat.Total); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return s, ErrTimedOut
		}
		if errors.Is(err, sql.ErrNoRows) {
			return s, ErrNoRows
		}

		return s, err
	}

	return s, nil
}

func (m *SnapshotModel) GetAllStats(ctx context.Context) ([]Snapshot, error) {
	query := `SELECT 
					strftime('%s', strftime('%Y-%m-%d', timestamp, 'unixepoch', 'localtime')) AS day_unix, SUM(sent), SUM(received), SUM(total)
					FROM snapshots
					GROUP BY day_unix
					ORDER BY day_unix DESC`

	timeout, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	rows, err := m.db.QueryContext(timeout, query)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, ErrTimedOut
		}

		return nil, err
	}

	var stats []Snapshot

	defer rows.Close()

	for rows.Next() {
		var s Snapshot

		if err = rows.Scan(&s.Timestamp, &s.Stat.Sent, &s.Stat.Received, &s.Stat.Total); err != nil {
			return nil, err
		}

		stats = append(stats, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}

func (m *SnapshotModel) GetMonthsInYear(ctx context.Context, year string) ([]string, error) {
	query := `
	SELECT DISTINCT strftime('%m', timestamp, 'unixepoch', 'localtime') AS month 
	FROM snapshots
	WHERE strftime('%Y', timestamp, 'unixepoch', 'localtime') = ?
	ORDER BY month DESC`

	timeout, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()
	rows, err := m.db.QueryContext(timeout, query, year)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, ErrTimedOut
		}

		return nil, err
	}

	defer rows.Close()

	var months []string

	for rows.Next() {
		var month string

		if err = rows.Scan(&month); err != nil {
			return nil, err
		}

		months = append(months, month)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return months, nil
}

func (m *SnapshotModel) GetStatByDate(ctx context.Context, date string) (DateStat, error) {
	query := `SELECT COUNT(*), SUM(sent), SUM(received), SUM(total)
	FROM (
		SELECT sent, received, total
		FROM snapshots
		WHERE strftime('%Y-%m-%d', timestamp, 'unixepoch', 'localtime') = ?
	)`

	timeout, cancel := context.WithTimeout(ctx, 3*time.Second)

	defer cancel()

	var s DateStat

	if err := m.db.QueryRowContext(timeout, query, date).Scan(&s.HoursMonitored, &s.Sent, &s.Received, &s.Total); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return s, ErrTimedOut
		}

		if errors.Is(err, sql.ErrNoRows) {
			return s, ErrNoRows
		}

		return s, err
	}

	return DateStat{
		HoursMonitored: s.HoursMonitored,
		Stat: Stat{
			Sent:     s.Stat.Sent,
			Received: s.Stat.Received,
			Total:    s.Stat.Total,
		},
	}, nil
}
