package add

import (
	"Motivation_reference/pkg/api/response"
	"Motivation_reference/pkg/logger"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator"
	"io"
	"net/http"
)

type Request struct {
	Name string `json:"name" validate:"required"`
}

type Response struct {
	response.Response
}

type addCategory interface {
	AddCategory(name string) (int64, error)
}

func New(logger logger.Logger, addCategory addCategory, w http.ResponseWriter, r *http.Request) {
	const op = "handlers.categories.New"
	w.Header().Set("Content-Type", "application/json")

	var req Request

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		if errors.Is(err, io.EOF) {
			logger.Error("request body is empty ", op)

			if err := json.NewEncoder(w).Encode(response.Error("empty request")); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			return
		}

		logger.Error("failed to decode request body", err, op)

		if err := json.NewEncoder(w).Encode(response.Error("failed to decode request")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	logger.Info("request body decoded", req)

	if err := validator.New().Struct(req); err != nil {
		var validateError validator.ValidationErrors
		errors.As(err, &validateError)

		logger.Error("invalid request", err)

		w.WriteHeader(http.StatusBadRequest)

		if err := json.NewEncoder(w).Encode(response.ValidationError(validateError)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	_, err = addCategory.AddCategory(req.Name)
	if err != nil {
		logger.Error("failed to add new category", err)

		w.WriteHeader(http.StatusBadRequest)

		if err := json.NewEncoder(w).Encode(response.Error("failed to add new addCategory")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	logger.Info("category added")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(Response{
		Response: response.OK(),
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
