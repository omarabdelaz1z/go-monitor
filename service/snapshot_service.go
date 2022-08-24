package service

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/omarabdelaz1z/go-monitor/model"
	"github.com/omarabdelaz1z/go-monitor/provider"
)

type SnapshotService interface {
	Create(*model.Snapshot) error
	GetStatsByMonth(string) ([]model.Snapshot, error)
	GetMonthStat(string) (model.MonthStat, error)
	GetAllStats() ([]model.Snapshot, error) // by day
}

type snapshotService struct {
	db *sql.DB
}

func NewSnapshotService(d *provider.DatabaseConnection) SnapshotService {
	return &snapshotService{db: d.GetDB()}
}

func (s *snapshotService) Create(snapshot *model.Snapshot) error {
	query := `INSERT INTO snapshots (timestamp, sent, received, total) VALUES (?, ?, ?, ?)`

	_, err := s.db.Exec(query,
		snapshot.Timestamp,
		snapshot.Stat.Sent,
		snapshot.Stat.Received,
		snapshot.Stat.Total,
	)

	if err != nil {
		return fmt.Errorf("failed to insert snapshot: %s", err)
	}

	return nil
}

func (s *snapshotService) GetStatsByMonth(month string) ([]model.Snapshot, error) {
	query := `SELECT 
			strftime('%s', strftime('%Y-%m-%d', timestamp, 'unixepoch')) AS day_unix,
			SUM(sent),
			SUM(received),
			SUM(total)
		FROM snapshots
		WHERE strftime('%m', timestamp, 'unixepoch') = ?
		GROUP BY day_unix
		ORDER BY day_unix DESC`

	rows, err := s.db.Query(query, month)

	if err != nil {
		return nil, fmt.Errorf("failed to query: %s", err)
	}

	var snapshots []model.Snapshot

	defer rows.Close()

	for rows.Next() {
		var snapshot model.Snapshot

		err := rows.Scan(&snapshot.Timestamp, &snapshot.Stat.Sent, &snapshot.Stat.Received, &snapshot.Stat.Total)

		if err != nil {
			return snapshots, fmt.Errorf("found no rows: %s", err)
		}

		snapshots = append(snapshots, snapshot)
	}

	return snapshots, nil
}

func (s *snapshotService) GetMonthStat(month string) (model.MonthStat, error) {
	var monthStat model.MonthStat

	query := `SELECT 
		strftime('%m', timestamp, 'unixepoch') AS month,
		SUM(sent), 
		SUM(received), 
		SUM(total) 
	FROM snapshots 
	WHERE month = ?
	GROUP BY month`

	row := s.db.QueryRow(query, month)

	if err := row.Scan(&monthStat.Month, &monthStat.Stat.Sent, &monthStat.Stat.Received, &monthStat.Stat.Total); err != nil {
		return monthStat, fmt.Errorf("found no rows: %s", err)
	}

	return monthStat, nil
}

func (s *snapshotService) GetAllStats() ([]model.Snapshot, error) {
	query := `SELECT 
	strftime('%s', strftime('%Y-%m-%d', timestamp, 'unixepoch')) AS day_unix,
	SUM(sent),
	SUM(received),
	SUM(total)
FROM snapshots
GROUP BY day_unix
ORDER BY day_unix DESC`

	rows, err := s.db.Query(query)

	if err != nil {
		return nil, fmt.Errorf("failed to query: %s", err)
	}

	var snapshots []model.Snapshot

	defer rows.Close()

	for rows.Next() {
		var snapshot model.Snapshot

		err := rows.Scan(&snapshot.Timestamp, &snapshot.Stat.Sent, &snapshot.Stat.Received, &snapshot.Stat.Total)

		if err != nil {
			return snapshots, fmt.Errorf("found no rows: %s", err)
		}

		snapshots = append(snapshots, snapshot)
	}

	return snapshots, nil
}
