package main

import (
	"aerospikedemo/internal/app"
	cfg "aerospikedemo/internal/app/config"
	"log"
	"time"
)

func main() {
	start := time.Now()
	log.Printf("Application is started at %s\n", start)

	config := cfg.NewConfig()
	config.StartTime = start

	handler := app.NewHandler(config)
	handler.Process()

	log.Printf("Total execution time %s", time.Since(start))

}
