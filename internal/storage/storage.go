package storage

import (
	"errors"
)

var (
	ErrCategoryExist = errors.New("this category already exists")
	ErrURLExists     = errors.New("this phrase exists")
)
