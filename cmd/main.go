package main

import (
	"Motivation_reference/internal/handlers"
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

	st, err := postgresql.New(cfg.Db.ConnString)
	if err != nil {
		logger.Fatal(err)
	}

	ph, err := st.UpgradePhrase(4, "newText")
	if err != nil {
		logger.Fatal(err)
	}

	fmt.Println(ph)

	http.HandleFunc("/api/v1/phrases", handlers.HandlerWithoutId)
	http.HandleFunc("/api/v1/phrases/{id}", handlers.HandlerWithId)

	logger.Infof("server started at %s:%s", cfg.Listen.BindIp, cfg.Listen.Port)

	http.ListenAndServe(fmt.Sprintf(":%s", cfg.Listen.Port), nil)
}
