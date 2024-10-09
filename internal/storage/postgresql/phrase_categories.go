package postgresql

import "fmt"

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

func (s *Storage) GetLink(phraseId, categoryId int64) error {
	const op = "storage.postgresql.GetLink"

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
	//TODO
	return nil
}

func (s *Storage) UpdateLink(phraseId, categoryId, newCategoryId int64) error {
	const op = "storage.postgresql.UpdateLink"
	//TODO
	return nil
}
