package storage

import (
	"database/sql"
	"errors"
)

var (
	ErrURLNotFound = errors.New("phrase not found")
	ErrURLExists   = errors.New("this phrase exists")
)

type Storage struct {
	DB *sql.DB
}
