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
}

type snapshotService struct {
	db *sql.DB
}

func NewSnapshotService(d *provider.Database) SnapshotService {
	return &snapshotService{db: d.GetDB()}
}

func (s *snapshotService) Create(stat *model.Snapshot) error {
	_, err := s.db.Exec(
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
