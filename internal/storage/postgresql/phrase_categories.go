package postgresql

import (
	"Motivation_reference/internal/storage"
	"errors"
	"fmt"
)

type Link struct {
	PhraseId   int64 `json:"phrase_id"`
	CategoryId int64 `json:"category_id"`
}

func (s *Storage) GetLink(phraseId, categoryId int64) (*Link, error) {
	const op = "storage.postgresql.GetLink"

	var link *Link

	err := s.db.QueryRow("SELECT * FROM phrase_categories WHERE phrase_id = $2 AND category_id = $2", phraseId, categoryId).Scan(&link)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return link, nil
}

func (s *Storage) GetLinkCategories(id int64) ([]Category, error) {
	const op = "storage.postgresql.GetLinkCategories"

	var categories []Category
	rows, err := s.db.Query("SELECT category_id FROM phrase_categories WHERE phrase_id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("%s: %d", op, err)
	}

	for rows.Next() {
		var categoryId int64
		if err = rows.Scan(&categoryId); err != nil {
			return nil, fmt.Errorf("%s: %d", op, err)
		}

		var category *Category
		category, err = s.GetCategory(categoryId)
		if err != nil {
			return nil, fmt.Errorf("%s: %d", op, err)
		}

		categories = append(categories, *category)
	}

	return categories, nil
}

func (s *Storage) GetLinkPhrases(id int64) ([]Phrase, error) {
	const op = "storage.postgresql.GetLinkCategories"

	var phrases []Phrase
	rows, err := s.db.Query("SELECT phrase_id FROM phrase_categories WHERE category_id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("%s: %d", op, err)
	}

	for rows.Next() {
		var phraseId int64
		if err = rows.Scan(&phraseId); err != nil {
			return nil, fmt.Errorf("%s: %d", op, err)
		}

		var phrase *Phrase
		phrase, err = s.GetPhrase(phraseId)
		if err != nil {
			return nil, fmt.Errorf("%s: %d", op, err)
		}

		phrases = append(phrases, *phrase)
	}

	return phrases, nil
}

func (s *Storage) AddLink(phraseId, categoryId int64) error {
	const op = "storage.postgresql.AddLink"

	stmt, err := s.db.Prepare("INSERT INTO phrase_categories VALUES ($1, $2)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(phraseId, categoryId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteLink(phraseId, categoryId int64) error {
	const op = "storage.postgresql.DeleteLink"

	stmt, err := s.db.Prepare("DELETE FROM phrase_categories WHERE phrase_id = $1 AND category_id = $2")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(phraseId, categoryId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdateLink(phraseId, categoryId int64, newCategoryName string) error {
	const op = "storage.postgresql.UpdateLink"

	newCategoryId, err := s.AddCategory(newCategoryName)
	if errors.Is(err, storage.ErrURLExists) {
		name, err := s.GetCategoryName(newCategoryName)
		if err != nil {
			return fmt.Errorf("%s: %d", op, err)
		}

		newCategoryId = name.Id
	} else if err != nil {
		return fmt.Errorf("%s: %d", op, err)
	}

	stmt, err := s.db.Prepare("UPDATE phrase_categories SET category_id = $1 WHERE phrase_id = $2 AND category_id = $3")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(newCategoryId, phraseId, categoryId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
