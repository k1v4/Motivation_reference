package get

import (
	"Motivation_reference/internal/storage/postgresql"
	"Motivation_reference/pkg/api/response"
	"Motivation_reference/pkg/logger"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type Response struct {
	response.Response
	Category postgresql.Category
}

type getCategory interface {
	GetCategory(id int64) (*postgresql.Category, error)
}

func New(logger logger.Logger, getCategory getCategory, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	arg := r.PathValue("id")
	id, err := strconv.Atoi(arg)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Errorf("failed to convert id. %s", err)

		if err := json.NewEncoder(w).Encode(response.Error("failed to convert id")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	logger.Info("got id")

	phrase, err := getCategory.GetCategory(int64(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		if errors.Is(err, sql.ErrNoRows) {
			logger.Errorf("no item with this id: %d. %s", id, err)

			if err := json.NewEncoder(w).Encode(response.Error("failed to get item")); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			return
		}

		logger.Error("failed to get item: ", err)

		if err := json.NewEncoder(w).Encode(response.Error("failed to get item")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(Response{
		Response: response.OK(),
		Category: postgresql.Category{
			Id:   phrase.Id,
			Name: phrase.Name,
		},
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
