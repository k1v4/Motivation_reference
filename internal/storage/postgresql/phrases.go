package postgresql

import (
	"Motivation_reference/internal/storage"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func New(connString string) (*storage.Storage, error) {
	const op = "storage.postgresql.New"

	db, err := sql.Open("pgx", connString)
	if err != nil {
		return &storage.Storage{}, fmt.Errorf("failed to open: %s: %d", op, err)
	}
	defer db.Close()

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS phrases(
		id INTEGER PRIMARY KEY,
		text TEXT NOT NULL UNIQUE);
	`)
	// CREATE INDEX IF NOT EXISTS idx_alias ON phrases(text);
	if err != nil {
		return nil, fmt.Errorf("%s: failed to prepare: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to execute: %w", op, err)
	}

	return &storage.Storage{DB: db}, nil
}
