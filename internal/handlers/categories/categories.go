package categories

import (
	"Motivation_reference/internal/handlers/categories/add"
	deleteCat "Motivation_reference/internal/handlers/categories/delete"
	"Motivation_reference/internal/handlers/categories/get"
	"Motivation_reference/internal/handlers/categories/getAll"
	"Motivation_reference/internal/handlers/categories/update"
	"Motivation_reference/internal/storage/postgresql"
	"Motivation_reference/pkg/logger"
	"net/http"
)

func HandlerWithoutId(logger *logger.Logger, storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getAll.New(*logger, storage, w, r)
		case http.MethodPost:
			add.New(*logger, storage, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	}
}

func HandlerWithId(logger *logger.Logger, storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			get.New(*logger, storage, w, r)
		case http.MethodDelete:
			deleteCat.New(*logger, storage, w, r)
		case http.MethodPatch:
			update.New(*logger, storage, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	}
}
