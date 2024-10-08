package postgresql

import (
	"Motivation_reference/internal/storage"
	"database/sql"
	"errors"
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

	// создание таблицы категорий
	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS categories (
    	id SERIAL PRIMARY KEY,
    	name VARCHAR(255) UNIQUE NOT NULL
	);
`)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to prepare: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to execute: %w", op, err)
	}

	// создание таблицы фраз
	stmt, err = db.Prepare(`
	CREATE TABLE IF NOT EXISTS phrases (
    	id SERIAL PRIMARY KEY,
    	text VARCHAR(255) UNIQUE NOT NULL
);`)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to prepare: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to execute: %w", op, err)
	}

	// создание таблицы для связи М:М
	stmt, err = db.Prepare(`
CREATE TABLE IF NOT EXISTS phrase_categories (
    phrase_id INTEGER REFERENCES phrases(id) ON DELETE CASCADE,
    category_id INTEGER REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY (phrase_id, category_id)
);`)

	if err != nil {
		return nil, fmt.Errorf("%s: failed to prepare: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to execute: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) AddPhrase(phraseText, nameCategory string) (int64, int64, error) {
	const op = "storage.postgresql.AddPhrase"

	var lastInsertedPhraseID int64
	err := s.db.QueryRow("INSERT INTO phrases(text) VALUES ($1) RETURNING id", phraseText).Scan(&lastInsertedPhraseID)
	if err != nil {
		return lastInsertedPhraseID, 0, fmt.Errorf("%s: %d", op, err)
	}

	lastInsertedCategoryID, err := s.AddCategory(nameCategory)
	if err != nil {
		if errors.Is(err, storage.ErrURLExists) {
			return lastInsertedPhraseID, lastInsertedCategoryID, nil
		}
		return lastInsertedPhraseID, lastInsertedCategoryID, fmt.Errorf("%s: %d", op, err)
	}

	return lastInsertedPhraseID, lastInsertedCategoryID, nil
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
