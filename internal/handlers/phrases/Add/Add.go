package Add

import (
	"Motivation_reference/pkg/logger"
	"net/http"
)

type AddPhrase interface {
	AddPhrase(phraseText string) (int64, error)
}

func New(logger logger.Logger, addPhrase AddPhrase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.phrases.Add"

	}

}
