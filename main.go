package main

import (
	cfg "aerospikedemo/internal/app/config"
	"aerospikedemo/internal/app/services/writer"
	"bufio"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	start := time.Now()

	//f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	//if err != nil {
	//	log.Fatalf("error opening file: %v", err)
	//}
	//defer f.Close()
	//
	//log.SetOutput(f)

	config, err := cfg.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	config.StartTime = start

	//readerService := reader.NewReaderService()
	writerService := writer.NewWriterService(config)

	log.Println(config.Host)
	log.Println(config.Port)

	file, err := os.Open("out.txt")
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	// Create a new scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Loop through the file and read each line
	for scanner.Scan() {
		line := strings.TrimSuffix(scanner.Text(), "\n")
		if err := writerService.UpdateCounter(line); err != nil {
			log.Fatal(err)
		}
	}

	// Check for errors during the scan
	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading file: %s", err)
	}

	writerService.GetAllResults()

	elapsed := time.Since(start)
	log.Printf("Execution time %s", elapsed)

}

// key -> IP
// namespace = DB name
// set = DB table
//
