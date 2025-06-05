package app

import (
	"aerospikedemo/internal/app/config"
	"aerospikedemo/internal/app/services/writer"
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type Handler struct {
	writerService *writer.Service

	config *config.Config
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		writerService: writer.NewWriterService(cfg),
		config:        cfg,
	}
}

func (h *Handler) consumeIPs(ch chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	batch := make([]string, 0, h.config.BatchSize)

	for val := range ch { // Receive values from the channel until it's closed
		batch = append(batch, val)
		if len(batch) == h.config.BatchSize {
			h.writerService.FlushBatch(batch)
			batch = make([]string, 0, h.config.BatchSize)
		}

	}
	if len(batch) > 0 {
		h.writerService.FlushBatch(batch)
	}
}

func (h *Handler) Process() {
	ips := make(chan string, h.config.Workers)

	// Count number of lines in input file
	file, err := os.Open(h.config.InputFilename)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	linesNum, err := lineCounter(file)
	if err != nil {
		log.Printf("Error counting lines number for %s\n", h.config.InputFilename)
	} else {
		log.Printf("Number of lines in %s is %v\n", h.config.InputFilename, linesNum)
	}

	// Open input file
	file, err = os.Open(h.config.InputFilename)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	log.Printf("Reading file %s is started\n", h.config.InputFilename)
	defer file.Close()

	var wg sync.WaitGroup
	for i := 0; i < h.config.Workers; i++ {
		wg.Add(1)
		go h.consumeIPs(ips, &wg)
	}
	log.Printf("Started %v workers\n", h.config.Workers)

	// Loop through the file and read each line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		val := strings.TrimSuffix(scanner.Text(), "\n")
		ips <- val
	}
	close(ips)
	wg.Wait()
	log.Printf("Reading file %s is completed in %s\n", h.config.InputFilename, time.Since(h.config.StartTime))

	// Check for errors during the scan
	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading file: %s", err)
	}

	if err := h.writerService.GetAllResults(); err != nil {
		log.Fatal(err)
	}
}
