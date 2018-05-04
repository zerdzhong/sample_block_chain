package main

import (
	. "./blockchain"
	"log"

	"fmt"
	"github.com/joho/godotenv"
	"strconv"
)

var blockChain BlockChain

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	blockChain = NewBlockChain()
	blockChain.AddBlock("Send 1 BTC to Ivan")
	blockChain.AddBlock("Send 2 more BTC to Ivan")

	for _, block := range blockChain.GetAllBlocks() {

		fmt.Printf("Prev. hash: %x\n", block.PrevHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}

	//networking.StartServer()
}
