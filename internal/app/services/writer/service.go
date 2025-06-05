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

const (
	IP_BIN    = "ip"
	COUNT_BIN = "count"
)

func NewWriterService(cfg *config.Config) *Service {
	s := new(Service)
	s.config = cfg

	var err error
	s.cli, err = aerospike.NewClient(cfg.Host, cfg.Port)
	if err != nil {
		log.Fatal(err)
	}

	s.writePolicy = aerospike.NewWritePolicy(0, 0)
	s.writePolicy.TotalTimeout = s.config.Timeout * time.Millisecond
	s.writePolicy.Expiration = uint32(s.config.TTL) // TTL 1h
	s.readPolicy = aerospike.NewPolicy()

	return s
}

func (s *Service) UpdateCounter(key string) error {
	k, err := aerospike.NewKey(s.config.Namespace, s.config.Set, key)
	if err != nil {
		log.Println("Error creating aerospike key")
		return err
	}

	// Attempt to get the record
	_, err = s.cli.Get(s.readPolicy, k)

	// Check if record exists
	if err != nil {
		if err.Matches(aerospike.ErrKeyNotFound.ResultCode) {
			// Create the record if it does not exist
			bins := aerospike.BinMap{}
			bins[IP_BIN] = key
			bins[COUNT_BIN] = 1
			err = s.cli.Put(s.writePolicy, k, bins)
			if err != nil {
				log.Println("Error creating aerospike record")
				return err
			}
		} else {
			log.Println("Error getting aerospike record")
			return err
		}
	} else {
		counterBin := aerospike.NewBin(COUNT_BIN, 1)
		err = s.cli.AddBins(s.writePolicy, k, counterBin)
		if err != nil {
			log.Println("Error incr bin to aerospike")
			return err
		}
	}

	return nil
}

func (s *Service) GetAllResults() error {
	// Create a scan policy
	policy := aerospike.NewScanPolicy()

	// Initialize the scan
	scan, asErr := s.cli.ScanAll(policy, s.config.Namespace, s.config.Set)
	if asErr != nil {
		log.Println("Error while scanning all results")
		return asErr
	}

	// If the file doesn't exist, create it, or append to the file
	f, err := os.Create(s.config.OutputFilename)
	if err != nil {
		log.Printf("Error creating result file: %s", s.config.OutputFilename)
		return err
	}
	defer f.Close()
	log.Printf("File %s is opened for writing\n", s.config.OutputFilename)

	// Process each record
	for rec := range scan.Results() {
		if rec.Err != nil {
			log.Printf("Error reading record: %v", rec.Err)
			continue
		}
		if _, err := f.Write([]byte(fmt.Sprintf("%v , count=%v\n", rec.Record.Bins[IP_BIN], rec.Record.Bins[COUNT_BIN]))); err != nil {
			log.Printf("Error processing record: %v", rec.Record)
			return err
		}
	}

	if _, err := f.Write([]byte(fmt.Sprintf("\nTotal time %s", time.Since(s.config.StartTime)))); err != nil {
		log.Printf("Error adding total time to result file")
		return err
	}

	log.Printf("Writing file %s is completed\n", s.config.OutputFilename)

	return nil
}
