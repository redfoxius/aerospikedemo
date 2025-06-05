package main

import (
	"aerospikedemo/internal/app"
	cfg "aerospikedemo/internal/app/config"
	"log"
	"time"
)

func main() {
	config := cfg.NewConfig()

	handler := app.NewHandler(config)
	handler.Process()

	log.Printf("Total execution time %s", time.Since(config.StartTime))
}
