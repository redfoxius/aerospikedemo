package config

import (
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
}

const (
	AEROSPIKE_SET     = "ip_counters"
	AEROSPIKE_SET_BIN = "count"
)

func NewConfig() (*Config, error) {

	aeroSpikeHost := os.Getenv("AEROSPIKE_HOST")
	aeroSpikePort, _ := strconv.Atoi(os.Getenv("AEROSPIKE_PORT"))
	aeroSpikeMode := os.Getenv("AEROSPIKE_MODE")
	aeroSpikeNamespace := os.Getenv("AEROSPIKE_NAMESPACE")

	return &Config{
		Host:      aeroSpikeHost,
		Port:      aeroSpikePort,
		Namespace: aeroSpikeNamespace,
		Set:       AEROSPIKE_SET,
		Bin:       AEROSPIKE_SET_BIN,
		Mode:      aeroSpikeMode,
	}, nil
}
