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

func (h *Handler) consumeIPs(ch chan string) {
	for val := range ch { // Receive values from the channel until it's closed
		if err := h.writerService.UpdateCounter(val); err != nil {
			log.Printf("Error updating counter: %s", err)
		}
	}
}

func (h *Handler) Process() {
	ips := make(chan string, 100)

	file, scanner := h.readerService.GetFileScanner()
	defer file.Close()

	for i := 0; i < h.config.Workers; i++ {
		go h.consumeIPs(ips)
	}
	log.Printf("Started %v workers\n", h.config.Workers)

	// Loop through the file and read each line
	for scanner.Scan() {
		val := strings.TrimSuffix(scanner.Text(), "\n")
		ips <- val
	}
	close(ips)
	log.Printf("Reading file %s is completed in %s\n", h.config.InputFilename, time.Since(h.config.StartTime))

	// Check for errors during the scan
	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading file: %s", err)
	}

	if err := h.writerService.GetAllResults(); err != nil {
		log.Fatal(err)
	}
}
