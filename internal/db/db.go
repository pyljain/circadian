package db

import (
	"database/sql"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type Db struct {
	conn *sql.DB
}

func New(relationDBToUse, connectionString string) (*Db, error) {
	conn, err := sql.Open(relationDBToUse, connectionString)
	if err != nil {
		return nil, err
	}

	db := &Db{
		conn: conn,
	}

	err = db.CreateHealthCheckResultsTable()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (d *Db) CreateHealthCheckResultsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS health_check_results (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		target_endpoint TEXT NOT NULL,
		http_method TEXT NOT NULL,
		callout_timestamp TIMESTAMP NOT NULL,
		response_code INTEGER NOT NULL,
		response TEXT,
		time_taken FLOAT NOT NULL
	)`

	_, err := d.conn.Exec(query)
	return err
}
