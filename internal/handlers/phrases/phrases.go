package phrases

import (
	"Motivation_reference/internal/handlers/phrases/Add"
	"Motivation_reference/internal/handlers/phrases/Get"
	"Motivation_reference/internal/storage/postgresql"
	"Motivation_reference/pkg/logger"
	"net/http"
)

func HandlerWithoutId(logger *logger.Logger, storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:

		case http.MethodPost:
			Add.New(*logger, storage, w, r)
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	}
}

func HandlerWithId(logger *logger.Logger, storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			Get.New(*logger, storage, w, r)
		case http.MethodDelete:
		case http.MethodPatch:
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	}
}
