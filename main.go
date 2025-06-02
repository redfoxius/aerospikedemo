package main

import (
	"fmt"
	"github.com/aerospike/aerospike-client-go/v8"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	aeroSpikeHost := os.Getenv("AEROSPIKE_HOST")
	aeroSpikePort, _ := strconv.Atoi(os.Getenv("AEROSPIKE_PORT"))
	//aeroSpikeMode := os.Getenv("AEROSPIKE_MODE")

	fmt.Println(aeroSpikeHost)
	fmt.Println(aeroSpikePort)

	client, err := aerospike.NewClient(aeroSpikeHost, aeroSpikePort)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	fmt.Println(aeroSpikeHost)
	fmt.Println(aeroSpikePort)
	fmt.Println(client.IsConnected())

}
