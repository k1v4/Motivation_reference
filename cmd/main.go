package main

import (
	"Motivation_reference/pkg/cfg"
	logger "Motivation_reference/pkg/logger"
)

func main() {
	logger := logger.GetLogger()

	cfg := cfg.GetConfig()

	_ = logger
	_ = cfg
}
