package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

func InitDB(path string) (*DB, error) {
	// Init SQLite
	sqlitedb, err := sql.Open("sqlite3", "./chime.db")
	if err != nil {
		return nil, err
	}

	db := &DB{sqlitedb}

	// Enable foreign key enforcement
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, err
	}

	// Init Models
	if err := db.initChatroomModel(); err != nil {
		return nil, err
	}
	if err := db.initMessagesModel(); err != nil {
		return nil, err
	}

	return db, nil
}

// Helpers
func scanRows[T any](rows *sql.Rows, scanFunc func(*T) error) ([]T, error) {
	defer rows.Close()
	var items []T
	for rows.Next() {
		var item T
		if err := scanFunc(&item); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
