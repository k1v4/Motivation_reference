package postgresql

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Phrase struct {
	Id   int64  `json:"id"`
	Text string `json:"text"`
}

type Storage struct {
	db *sql.DB
}

func New(connString string) (*Storage, error) {
	const op = "storage.postgresql.New"

	db, err := sql.Open("pgx", connString)
	if err != nil {
		return &Storage{}, fmt.Errorf("failed to open: %s: %d", op, err)
	}
	//defer db.Close()

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS phrases(
		id SERIAL PRIMARY KEY,
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

	stmt, err = db.Prepare(`
	CREATE INDEX IF NOT EXISTS idx_text ON phrases(text);
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: failed to prepare: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to execute: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) AddPhrase(phraseText string) (int64, error) {
	const op = "storage.postgresql.AddPhrase"

	var lastInsertedID int64
	err := s.db.QueryRow("INSERT INTO phrases(text) VALUES ($1) RETURNING id", phraseText).Scan(&lastInsertedID)
	if err != nil {
		return lastInsertedID, fmt.Errorf("%s: %d", op, err)
	}

	return lastInsertedID, nil
}

func (s *Storage) GetPhrase(id int64) (*Phrase, error) {
	const op = "storage.postgresql.GetPhrase"

	var phrase Phrase
	err := s.db.QueryRow("SELECT id, text FROM phrases WHERE id = $1", id).Scan(&phrase.Id, &phrase.Text)
	if err != nil {
		return nil, fmt.Errorf("%s: %d", op, err)
	}

	return &phrase, nil
}

func (s *Storage) GetPhrases() ([]Phrase, error) {
	const op = "storage.postgresql.GetPhrases"

	var phrases []Phrase
	rows, err := s.db.Query("SELECT * FROM phrases")
	if err != nil {
		return nil, fmt.Errorf("%s: %d", op, err)
	}

	for rows.Next() {
		var phrase Phrase
		if err := rows.Scan(&phrase.Id, &phrase.Text); err != nil {
			return nil, fmt.Errorf("%s: %d", op, err)
		}

		phrases = append(phrases, phrase)
	}

	return phrases, nil
}

func (s *Storage) DeletePhrase(id int64) error {
	const op = "storage.postgresql.DeletePhrase"

	stmt, err := s.db.Prepare("DELETE FROM phrases WHERE id = $1")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UpgradePhrase(id int64, newText string) (*Phrase, error) {
	const op = "storage.postgresql.UpgradePhrase"

	stmt, err := s.db.Prepare("UPDATE phrases SET text=$1 WHERE id=$2")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(newText, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	phrase, err := s.GetPhrase(id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return phrase, nil
}
