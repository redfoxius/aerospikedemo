package writer

import (
	"aerospikedemo/internal/app/config"
	"fmt"
	"github.com/aerospike/aerospike-client-go/v8"
	"log"
	"os"
	"time"
)

type Service struct {
	config *config.Config

	cli         *aerospike.Client
	writePolicy *aerospike.WritePolicy
	readPolicy  *aerospike.BasePolicy
}

func NewWriterService(cfg *config.Config) *Service {
	service := new(Service)
	service.config = cfg

	var err error
	service.cli, err = aerospike.NewClient(cfg.Host, cfg.Port)
	if err != nil {
		log.Fatal(err)
	}

	service.writePolicy = aerospike.NewWritePolicy(0, 0)
	service.writePolicy.TotalTimeout = 5000 * time.Millisecond
	service.readPolicy = aerospike.NewPolicy()

	return service
}

func (s *Service) UpdateCounter(key string) error {
	k, err := aerospike.NewKey(s.config.Namespace, s.config.Set, key)
	if err != nil {
		log.Fatal(err)
	}

	counterBin := aerospike.NewBin(s.config.Bin, 1)
	err = s.cli.AddBins(s.writePolicy, k, counterBin)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (s *Service) GetAllResults() {
	// Create a scan policy
	policy := aerospike.NewScanPolicy()

	// Initialize the scan
	scan, asErr := s.cli.ScanAll(policy, s.config.Namespace, s.config.Set)
	if asErr != nil {
		log.Fatal(asErr)
	}

	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(s.config.OutputFilename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.Printf("File %s is opened for writing\n", s.config.OutputFilename)

	// Process each record
	for rec := range scan.Results() {
		if rec.Err != nil {
			log.Printf("Error reading record: %v", rec.Err)
			continue
		}
		if _, err := f.Write([]byte(fmt.Sprintf("%s , count=%s", rec.Record.Key.String(), rec.Record.Bins[s.config.Bin].(int)))); err != nil {
			log.Fatal(err)
		}
	}

	if _, err := f.Write([]byte(fmt.Sprintf("Total time %s", time.Since(s.config.StartTime)))); err != nil {
		log.Fatal(err)
	}

	log.Printf("Writing file %s is completed\n", s.config.OutputFilename)
}
