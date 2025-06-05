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

type BatchRecord struct {
	key  *aerospike.Key
	bins aerospike.BinMap
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

func (s *Service) FlushBatch(batch []string) error {
	batchRecords := make([]aerospike.BatchRecordIfc, 0, s.config.BatchSize)
	for _, k := range batch {
		key, err := aerospike.NewKey(s.config.Namespace, s.config.Set, k)
		if err != nil {
			log.Printf("Error creating aerospike key: %s\n", k)
			continue
		}
		record := aerospike.NewBatchWrite(
			nil,
			key,
			aerospike.PutOp(aerospike.NewBin(IP_BIN, k)),
			aerospike.AddOp(aerospike.NewBin(COUNT_BIN, 1)),
		)
		batchRecords = append(batchRecords, record)
	}

	err := s.cli.BatchOperate(nil, batchRecords)
	if err != nil {
		log.Printf("Error flushing batch records: %s", err)
		return err
	}

	return nil
}

func (s *Service) GetAllResults() error {
	totalRecords := uint32(0)

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
		} else {
			totalRecords += uint32(rec.Record.Bins[COUNT_BIN].(int))
		}
	}

	if _, err := f.Write([]byte(fmt.Sprintf("\nTotal time %s", time.Since(s.config.StartTime)))); err != nil {
		log.Printf("Error adding total time to result file")
		return err
	}
	if _, err := f.Write([]byte(fmt.Sprintf("\nTotal count %v", totalRecords))); err != nil {
		log.Printf("Error adding total count to result file")
		return err
	}

	log.Printf("Writing file %s is completed\n", s.config.OutputFilename)

	return nil
}
