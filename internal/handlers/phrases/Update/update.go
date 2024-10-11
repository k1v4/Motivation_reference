package Update

import (
	"Motivation_reference/internal/storage/postgresql"
	"Motivation_reference/pkg/api/response"
	"Motivation_reference/pkg/logger"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
)

type Request struct {
	Text            string `json:"text" validate:"required"`
	CategoryName    string `json:"category" validate:"required"`
	NewCategoryName string `json:"new_category" validate:"required"`
}

type Response struct {
	response.Response
	Phrase postgresql.Phrase
}

type updatePhrase interface {
	UpgradePhrase(id int64, newText, category, newCategory string) (*postgresql.Phrase, error)
}

func New(logger logger.Logger, updatePhrase updatePhrase, w http.ResponseWriter, r *http.Request) {
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

	var req Request

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		if errors.Is(err, io.EOF) {
			logger.Error("request body is empty")

			if err := json.NewEncoder(w).Encode(response.Error("empty request")); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			return
		}

		logger.Error("failed to decode request body", err)

		if err := json.NewEncoder(w).Encode(response.Error("failed to decode request")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	logger.Info("request body decoded")

	phrase, err := updatePhrase.UpgradePhrase(int64(id), req.Text, req.CategoryName, req.NewCategoryName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Errorf("failed to delete phrase. %s", err)

		if err := json.NewEncoder(w).Encode(response.Error("failed to delete phrase")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(Response{
		Response: response.OK(),
		Phrase:   *phrase,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
