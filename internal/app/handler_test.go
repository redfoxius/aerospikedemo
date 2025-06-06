package app

import (
	"aerospikedemo/internal/app/config"
	"aerospikedemo/internal/app/services/writer"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewHandler(t *testing.T) {
	cfg := getConfig()

	exp := &Handler{
		writerService: writer.NewWriterService(cfg),
		config:        cfg,
	}
	act := NewHandler(cfg)

	assert.Equal(t, exp, act)
}

func TestHandler_Process(t *testing.T) {
	handler := NewHandler(getConfig())

	err := handler.Process()

	assert.NotEqual(t, nil, err)
}

func getConfig() *config.Config {
	return &config.Config{
		Host:           "localhost",
		Port:           3000,
		Namespace:      "test",
		Set:            "ip_count",
		Timeout:        time.Duration(500),
		TTL:            300,
		Mode:           "sync",
		InputFilename:  "ip.txt",
		OutputFilename: "result.txt",
		Workers:        2,
		BatchSize:      10,
		StartTime:      time.Now(),
	}
}
