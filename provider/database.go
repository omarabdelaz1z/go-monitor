package provider

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	maxIdleConns    = 5
	maxOpenConns    = 10
	maxConnIdleTime = 2 * time.Hour
)

type Database struct {
	db *sql.DB
}

func NewConnection(driver, dsn string) (*Database, error) {
	db, err := sql.Open(driver, dsn)

	if err != nil {
		return nil, fmt.Errorf("failed to open database: %s", err)
	}

	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxIdleTime(maxConnIdleTime)

	return &Database{db}, nil
}

func (d *Database) GetDB() *sql.DB {
	return d.db
}
