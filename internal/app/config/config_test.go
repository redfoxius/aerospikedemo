package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestNewConfig(t *testing.T) {
	aeroSpikeHost := "localhost"
	aeroSpikePort := 6379
	aeroSpikeNamespace := "test"
	aeroSpikeSetName := "test"
	aeroSpikeTimeout := "500"
	aeroSpikeTTL := 500
	workersNum := 100
	batchSize := 1000
	startTime := time.Now()

	t.Setenv("AEROSPIKE_HOST", aeroSpikeHost)
	t.Setenv("AEROSPIKE_PORT", strconv.Itoa(aeroSpikePort))
	t.Setenv("AEROSPIKE_NAMESPACE", aeroSpikeNamespace)
	t.Setenv("AEROSPIKE_SET", aeroSpikeSetName)
	t.Setenv("AEROSPIKE_TIMEOUT", aeroSpikeTimeout)
	t.Setenv("AEROSPIKE_TTL", strconv.Itoa(aeroSpikeTTL))
	t.Setenv("WORKERS", strconv.Itoa(workersNum))
	t.Setenv("BATCH_SIZE", strconv.Itoa(batchSize))

	d, _ := time.ParseDuration(os.Getenv("AEROSPIKE_TIMEOUT"))
	exp := &Config{
		Host:           aeroSpikeHost,
		Port:           aeroSpikePort,
		Namespace:      aeroSpikeNamespace,
		Set:            aeroSpikeSetName,
		Timeout:        d,
		TTL:            aeroSpikeTTL,
		Mode:           "sync",
		InputFilename:  "ip.txt",
		OutputFilename: "result.txt",
		Workers:        workersNum,
		BatchSize:      batchSize,
		StartTime:      startTime,
	}

	act := NewConfig()
	act.StartTime = startTime

	assert.Equal(t, exp, act)
}
