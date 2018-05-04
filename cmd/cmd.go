package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
	"samplechain/blockchain"
	"strconv"
)

const (
	AddBlockCMD   = "addBlock"
	PrintChainCMD = "printChain"
)

type CMD struct {
	bc *blockchain.Blockchain
}

func NewCMD() *CMD {
	bc := blockchain.NewBlockChain()
	cmd := CMD{bc}
	return &cmd
}

func (cmd *CMD) Run() {

	err := cmd.validateArgs()

	addBlockCmd := flag.NewFlagSet(AddBlockCMD, flag.ExitOnError)
	printChainCmd := flag.NewFlagSet(PrintChainCMD, flag.ExitOnError)

	addBlockData := addBlockCmd.String("data", "", "Block data")

	switch os.Args[1] {
	case AddBlockCMD:
		err = addBlockCmd.Parse(os.Args[2:])
	case PrintChainCMD:
		err = printChainCmd.Parse(os.Args[2:])

	default:
		cmd.printUsage()
	}

	if nil != err {
		log.Fatal(err.Error())
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cmd.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cmd.printChain()
	}
}

func (cmd *CMD) validateArgs() error {
	if len(os.Args) < 2 {
		cmd.printUsage()
		os.Exit(1)
	}

	return nil
}

func (cmd *CMD) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  addBlock -data BLOCK_DATA - add a block to the bc")
	fmt.Println("  printChain - print all the blocks of the bc")
}

func (cmd *CMD) addBlock(data string) {
	err := cmd.bc.AddBlock(data)
	if nil != err {
		log.Fatal(err.Error())
	}
	fmt.Println("AddBlock Success!")
}

func (cmd *CMD) printChain() {
	bci := cmd.bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf(block.Description())
		pow := blockchain.NewProofOfWork(*block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break
		}
	}

}

func (cmd *CMD) Close() {
	cmd.bc.CloseDB()
}
