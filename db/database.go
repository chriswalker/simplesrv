// Package db handles all database operations (currently against
// SQLite).
package db

import (
	"database/sql"
	_ "embed"
	"time"

	"github.com/chriswalker/simplesrv/model"
	_ "modernc.org/sqlite"
)

//go:embed sql/schema.sql
var schema string

// DB holds the created SQLite DB structure.
type DB struct {
	db *sql.DB
}

// New opens the supplied SQLite database, and applies the base
// schema if needed.
func New(name string) (*DB, error) {
	db, err := sql.Open("sqlite", name)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Check if any tables exist. If they don't we'll run the schema
	// file to create & populate with some sample items.
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		if _, err := db.Exec(schema); err != nil {
			return nil, err
		}
	}
	db.SetMaxOpenConns(1)

	return &DB{db: db}, nil
}

// GetItems returns all items.
func (d *DB) GetItems() ([]model.Item, error) {
	rows, err := d.db.Query("SELECT * FROM items")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.Item{}
	for rows.Next() {
		i := model.Item{}
		var created int64
		err := rows.Scan(&i.ID, &i.Title, &i.Description, &created)
		if err != nil {
			return nil, err
		}
		i.CreationDate = time.Unix(created, 0).UTC()

		items = append(items, i)
	}

	return items, nil
}
