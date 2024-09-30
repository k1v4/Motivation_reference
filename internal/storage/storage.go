package storage

import (
	"errors"
)

var (
	ErrURLNotFound = errors.New("phrase not found")
	ErrURLExists   = errors.New("this phrase exists")
)
