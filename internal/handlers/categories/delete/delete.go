package delete

import (
	"Motivation_reference/pkg/api/response"
	"Motivation_reference/pkg/logger"
	"encoding/json"
	"net/http"
	"strconv"
)

type Response struct {
	response.Response
}

type deleteCategory interface {
	DeleteCategory(id int64) error
}

func New(logger logger.Logger, deleteCategory deleteCategory, w http.ResponseWriter, r *http.Request) {
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

	err = deleteCategory.DeleteCategory(int64(id))
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
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
