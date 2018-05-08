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
	createBlockChainCMDName = "createblockchain"
	getBalanceCMDName       = "getbalance"
	sendCMDName             = "send"
	printChainCMDName       = "printchain"
)

type CMD struct {
}

func (cmd *CMD) Run() {

	err := cmd.validateArgs()

	createBlockChainCMD := flag.NewFlagSet(createBlockChainCMDName, flag.ExitOnError)
	printChainCMD := flag.NewFlagSet(printChainCMDName, flag.ExitOnError)
	getBalanceCMD := flag.NewFlagSet(getBalanceCMDName, flag.ExitOnError)
	sendCMD := flag.NewFlagSet(sendCMDName, flag.ExitOnError)

	createBlockchainAddrData := createBlockChainCMD.String("address", "", "The address to send genesis block reward to")
	getBalanceAddrData := getBalanceCMD.String("address", "", "The address to get balance for")
	sendFrom := sendCMD.String("from", "", "Source wallet address")
	sendTo := sendCMD.String("to", "", "Destination wallet address")
	sendAmount := sendCMD.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case createBlockChainCMDName:
		err = createBlockChainCMD.Parse(os.Args[2:])
	case printChainCMDName:
		err = printChainCMD.Parse(os.Args[2:])
	case getBalanceCMDName:
		err = getBalanceCMD.Parse(os.Args[2:])
	case sendCMDName:
		err = sendCMD.Parse(os.Args[2:])

	default:
		cmd.printUsage()
	}

	if nil != err {
		log.Fatal(err.Error())
	}

	if createBlockChainCMD.Parsed() {
		if *createBlockchainAddrData == "" {
			createBlockChainCMD.Usage()
			os.Exit(1)
		}
		cmd.createBlockchain(*createBlockchainAddrData)
	}

	if getBalanceCMD.Parsed() {
		if *getBalanceAddrData == "" {
			createBlockChainCMD.Usage()
			os.Exit(1)
		}
		cmd.getBalance(*getBalanceAddrData)
	}

	if sendCMD.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCMD.Usage()
			os.Exit(1)
		}
		cmd.send(*sendFrom, *sendTo, *sendAmount)
	}

	if printChainCMD.Parsed() {
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

func (cmd *CMD) createBlockchain(address string) {
	bc := blockchain.NewBlockchain(address)
	bc.CloseDB()
	fmt.Println("Done!")
}

func (cmd *CMD) getBalance(address string) {
	bc := blockchain.NewBlockchain(address)
	defer bc.CloseDB()

	balance := 0
	UTXOs := bc.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}

func (cmd *CMD) send(from, to string, amount int) {
	bc := blockchain.NewBlockchain(from)
	defer bc.CloseDB()

	tx := blockchain.NewUTXOTransaction(from, to, amount, bc)
	bc.MineBlock([]*blockchain.Transaction{tx})
	fmt.Println("Success!")
}

func (cmd *CMD) printChain() {
	// TODO: Fix this
	bc := blockchain.NewBlockchain("")
	defer bc.CloseDB()

	bci := bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf(block.Description())
		pow := blockchain.NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break
		}
	}

}
