package reader

import (
	"aerospikedemo/internal/app/config"
	"bufio"
	"log"
	"os"
)

type Service struct {
	config *config.Config
}

func NewReaderService(cfg *config.Config) *Service {
	return &Service{cfg}
}

func (s *Service) GetFileScanner() (*os.File, *bufio.Scanner) {
	// Open input file
	file, err := os.Open(s.config.InputFilename)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	log.Printf("Reading file %s is started\n", s.config.InputFilename)

	// Create a new scanner to read the file line by line
	return file, bufio.NewScanner(file)
}
