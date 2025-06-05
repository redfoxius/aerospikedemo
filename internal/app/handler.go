package app

import (
	"aerospikedemo/internal/app/config"
	"aerospikedemo/internal/app/services/reader"
	"aerospikedemo/internal/app/services/writer"
	"log"
	"strings"
	"time"
)

type Handler struct {
	writerService *writer.Service
	readerService *reader.Service

	config *config.Config
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		writerService: writer.NewWriterService(cfg),
		readerService: reader.NewReaderService(cfg),
		config:        cfg,
	}
}

func (f *Handler) Process() {
	file, scanner := f.readerService.GetFileScanner()
	defer file.Close()

	// Loop through the file and read each line
	for scanner.Scan() {
		line := strings.TrimSuffix(scanner.Text(), "\n")
		if err := f.writerService.UpdateCounter(line); err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("Reading file %s is completed in %s\n", f.config.InputFilename, time.Since(f.config.StartTime))

	// Check for errors during the scan
	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading file: %s", err)
	}

	f.writerService.GetAllResults()
}
