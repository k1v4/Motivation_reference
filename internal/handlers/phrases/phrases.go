package phrases

import (
	"Motivation_reference/internal/handlers/phrases/Add"
	"Motivation_reference/internal/handlers/phrases/Delete"
	"Motivation_reference/internal/handlers/phrases/Get"
	"Motivation_reference/internal/handlers/phrases/GetAll"
	"Motivation_reference/internal/handlers/phrases/Update"
	"Motivation_reference/internal/storage/postgresql"
	"Motivation_reference/pkg/logger"
	"net/http"
)

func HandlerWithoutId(logger *logger.Logger, storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			GetAll.New(*logger, storage, w, r)
		case http.MethodPost:
			Add.New(*logger, storage, w, r)
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
			Get.New(*logger, storage, w, r)
		case http.MethodDelete:
			Delete.New(*logger, storage, w, r)
		case http.MethodPatch:
			Update.New(*logger, storage, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	}
}
