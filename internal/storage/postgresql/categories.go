package postgresql

import (
	"Motivation_reference/internal/storage"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
)

type Category struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func (s *Storage) AddCategory(name string) (int64, error) {
	const op = "storage.postgresql.AddCategory"

	var lastInsertedID int64
	err := s.db.QueryRow(`
		INSERT INTO categories(name) 
		VALUES ($1) 
		ON CONFLICT (name) DO UPDATE
		SET name = EXCLUDED.name 
		RETURNING id
	`, name).Scan(&lastInsertedID)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return lastInsertedID, storage.ErrCategoryExist
			}
		}

		return lastInsertedID, fmt.Errorf("%s: %d", op, err)
	}

	return lastInsertedID, nil
}

func (s *Storage) GetCategory(id int64) (*Category, error) {
	const op = "storage.postgresql.GetCategory"

	var category Category
	err := s.db.QueryRow("SELECT id, name FROM categories WHERE id = $1", id).Scan(&category.Id, &category.Name)
	if err != nil {
		return nil, fmt.Errorf("%s: %d", op, err)
	}

	return &category, nil
}

func (s *Storage) GetCategories() ([]Category, error) {
	const op = "storage.postgresql.GetCategories"

	var categories []Category
	rows, err := s.db.Query("SELECT * FROM categories")
	if err != nil {
		return nil, fmt.Errorf("%s: %d", op, err)
	}

	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.Id, &category.Name); err != nil {
			return nil, fmt.Errorf("%s: %d", op, err)
		}

		categories = append(categories, category)
	}

	return categories, nil
}

func (s *Storage) DeleteCategory(id int64) error {
	const op = "storage.postgresql.DeleteCategory"

	stmt, err := s.db.Prepare("DELETE FROM categories WHERE id = $1")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UpgradeCategory(id int64, newName string) (*Category, error) {
	const op = "storage.postgresql.UpgradeCategory"

	stmt, err := s.db.Prepare("UPDATE categories SET name=$1 WHERE id=$2")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(newName, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	category, err := s.GetCategory(id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return category, nil
}
