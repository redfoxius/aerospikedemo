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
	Config *config.Config

	cli         *aerospike.Client
	writePolicy *aerospike.WritePolicy
	readPolicy  *aerospike.BasePolicy
}

func NewWriterService(cfg *config.Config) *Service {
	service := new(Service)
	service.Config = cfg

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
	k, err := aerospike.NewKey(s.Config.Namespace, s.Config.Set, key)
	if err != nil {
		log.Fatal(err)
	}

	counterBin := aerospike.NewBin(s.Config.Bin, 1)
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
	scan, err := s.cli.ScanAll(policy, s.Config.Namespace, s.Config.Set)
	if err != nil {
		log.Fatal(err)
	}

	// If the file doesn't exist, create it, or append to the file
	f, er := os.OpenFile("app/result.txt", os.O_CREATE|os.O_WRONLY, 0644)
	if er != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Process each record
	for rec := range scan.Results() {
		if rec.Err != nil {
			log.Printf("Error reading record: %v", rec.Err)
			continue
		}
		if _, err := f.Write([]byte(fmt.Sprintf("%s , count=%s", rec.Record.Key.String(), rec.Record.Bins[s.Config.Bin].(int)))); err != nil {
			log.Fatal(err)
		}
	}

	elapsed := time.Since(s.Config.StartTime)
	if _, err := f.Write([]byte(fmt.Sprintf("Total time %s", elapsed))); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Scan complete")
}
