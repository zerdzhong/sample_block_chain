package main

import (
	"log"
	"samplechain/cmd"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	cmd := cmd.CMD{}
	cmd.Run()
}
