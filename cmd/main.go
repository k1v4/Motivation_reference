package main

import (
	"Motivation_reference/internal/handlers/categories"
	"Motivation_reference/internal/handlers/phrases"
	"Motivation_reference/internal/storage/postgresql"
	"Motivation_reference/pkg/cfg"
	logger "Motivation_reference/pkg/logger"
	"fmt"
	"net/http"
	"time"
)

func main() {
	logger := logger.GetLogger()

	cfg := cfg.GetConfig()

	// GET "/api/v1/phrases"
	// GET(one) "/api/v1/phrases/{id}"
	// POST "/api/v1/phrases"
	// DELETE "/api/v1/phrases/{id}"
	// PATCH "/api/v1/phrases/{id}"

	// GET /api/v1/categories
	// GET(one) "/api/v1/categories/{id}"
	// POST "/api/v1/categories"
	// DELETE "/api/v1/categories/{id}"
	// PATCH "/api/v1/categories/{id}"

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.Postgresql.Username,
		cfg.Postgresql.Password,
		cfg.Postgresql.Host,
		cfg.Postgresql.Port,
		cfg.Postgresql.Database)
	var storage *postgresql.Storage
	var err error

	for i := 0; i < 3; i++ {
		storage, err = postgresql.New(connString)
		if err != nil {
			logger.Error(err)
			time.Sleep(2 * time.Second)
		}
	}

	logger.Info("DB connected")

	http.HandleFunc("/api/v1/phrases", phrases.HandlerWithoutId(logger, storage))
	http.HandleFunc("/api/v1/phrases/{id}", phrases.HandlerWithId(logger, storage))
	http.HandleFunc("/api/v1/categories", categories.HandlerWithoutId(logger, storage))
	http.HandleFunc("/api/v1/categories/{id}", categories.HandlerWithId(logger, storage))

	logger.Infof("server started at %s:%s", cfg.Listen.BindIp, cfg.Listen.Port)

	http.ListenAndServe(fmt.Sprintf(":%s", cfg.Listen.Port), nil)
}
