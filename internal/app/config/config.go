package config

import (
	"flag"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	// AeroSpike config params
	Host      string
	Port      int
	Namespace string
	Set       string
	Bin       string

	Mode string

	StartTime time.Time

	InputFilename  string
	OutputFilename string
}

const (
	AEROSPIKE_SET     = "ip_counters"
	AEROSPIKE_SET_BIN = "count"
)

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	aeroSpikeHost := os.Getenv("AEROSPIKE_HOST")
	aeroSpikePort, _ := strconv.Atoi(os.Getenv("AEROSPIKE_PORT"))
	aeroSpikeNamespace := os.Getenv("AEROSPIKE_NAMESPACE")

	inputPtr := flag.String("in", "ip.txt", "Filename/path for input file, a string.")
	outputPtr := flag.String("out", "result.txt", "Filename/path for output file, a string.")
	modePtr := flag.String("mode", "sync", "Mode for working with AeroSpike (sync/async), a string.")
	flag.Parse()

	if *modePtr != "async" && *modePtr != "sync" {
		log.Fatal("Mode must be sync or async")
	}

	return &Config{
		Host:           aeroSpikeHost,
		Port:           aeroSpikePort,
		Namespace:      aeroSpikeNamespace,
		Set:            AEROSPIKE_SET,
		Bin:            AEROSPIKE_SET_BIN,
		Mode:           *modePtr,
		InputFilename:  *inputPtr,
		OutputFilename: *outputPtr,
	}
}
