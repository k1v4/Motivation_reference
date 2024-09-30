package main

import (
	"Motivation_reference/internal/handlers/phrases"
	"Motivation_reference/internal/storage/postgresql"
	"Motivation_reference/pkg/cfg"
	logger "Motivation_reference/pkg/logger"
	"fmt"
	"net/http"
)

func main() {
	logger := logger.GetLogger()

	cfg := cfg.GetConfig()

	_ = logger
	_ = cfg

	// GET "/api/v1/phrases"
	// GET(one) "/api/v1/phrases/{id}"
	// POST "/api/v1/phrases"
	// DELETE "/api/v1/phrases/{id}"
	// PATCH "/api/v1/phrases/{id}"

	storage, err := postgresql.New(cfg.Db.ConnString)
	if err != nil {
		logger.Fatal(err)
	}
	http.HandleFunc("/api/v1/phrases", phrases.HandlerWithoutId(logger, storage))
	http.HandleFunc("/api/v1/phrases/{id}", phrases.HandlerWithId(logger, storage))

	logger.Infof("server started at %s:%s", cfg.Listen.BindIp, cfg.Listen.Port)

	http.ListenAndServe(fmt.Sprintf(":%s", cfg.Listen.Port), nil)
}
