package main

import (
	"github.com/joho/godotenv"
	"log"
	"samplechain/blockchain"
	"samplechain/cmd"
	"samplechain/http"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	go http.Run()

	cmd := cmd.CMD{}
	cmd.Run()

	defer blockchain.SampleChain.CloseDB()
}
