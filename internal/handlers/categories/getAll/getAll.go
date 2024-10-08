package getAll

import (
	"Motivation_reference/internal/storage/postgresql"
	"Motivation_reference/pkg/api/response"
	"Motivation_reference/pkg/logger"
	"encoding/json"
	"net/http"
)

type Response struct {
	response.Response
	Category []postgresql.Category
}

type getAllCategories interface {
	GetCategories() ([]postgresql.Category, error)
}

func New(logger logger.Logger, getAll getAllCategories, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	categories, err := getAll.GetCategories()
	if err != nil {
		logger.Errorf("failed to get all phrases: %s", err)

		w.WriteHeader(http.StatusBadRequest)

		if err := json.NewEncoder(w).Encode(response.Error("failed to get all phrases")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(Response{
		Response: response.OK(),
		Category: categories,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
