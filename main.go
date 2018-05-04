package main

import (
	"log"
	. "samplechain/cmd"

	"github.com/joho/godotenv"
	"samplechain/networking"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	cmd := NewCMD()

	defer cmd.Close()
	cmd.Run()

	networking.StartServer()
}
