package Get

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
	Phrase postgresql.Phrase
}

type getPhrase interface {
	GetPhrase(id int64) (*postgresql.Phrase, error)
}

func New(logger logger.Logger, getPhrase getPhrase, w http.ResponseWriter, r *http.Request) {
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

	phrase, err := getPhrase.GetPhrase(int64(id))
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
		Phrase: postgresql.Phrase{
			Id:   phrase.Id,
			Text: phrase.Text,
		},
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
