package main

import (
	. "./blockchain"
	"log"

	"github.com/joho/godotenv"

	"samplechain/networking"
)

var blockChain BlockChain

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	blockChain = NewBlockChain()

	networking.StartServer()
}
