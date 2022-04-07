package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"github.com/jaczerob/madamchuckle/internal/static"
)

type Database struct {
	db *sql.DB
}

func New() (d *Database, err error) {
	db, err := sql.Open("sqlite3", static.DBPath)
	if err != nil {
		return
	}

	stmt := `
	CREATE TABLE IF NOT EXISTS commands (
		id 		VARCHAR(64) NOT NULL,
		command VARCHAR(12) NOT NULL
	);

	CREATE TABLE IF NOT EXISTS events (
		event_id 	INTEGER 	CHECK(event_id in (0, 1, 2, 3, 4, 5)) NOT NULL,
		message_id 	VARCHAR(18) NOT NULL,
		channel_id 	VARCHAR(18) NOT NULL
	);
	`

	_, err = db.Exec(stmt)
	if err != nil {
		return
	}

	d = &Database{db: db}
	return
}
