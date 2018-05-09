package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
	"samplechain/p2p"
)

type CMD struct {
}

func (cmd *CMD) Run() {

	err := cmd.validateArgs()

	listenF := flag.Int("l", 0, "wait for incoming connections")
	target := flag.String("d", "", "target peer to dial")
	secio := flag.Bool("secio", false, "enable secureIO")
	seed := flag.Int64("seed", 0, "set random seed for id generation")
	flag.Parse()

	if *listenF == 0 {
		log.Fatal("Please provide a port to bind on with -l")
	}

	// Make a host that listens on the given multiaddress
	ha, err := p2p.MakeBasicHost(*listenF, *secio, *seed)
	if err != nil {
		log.Fatal(err)
	}

	flag.Parse()

	if *target == "" {
		select {} // hang forever
	} else {
		p2p.DialToTarget(ha, *target)
		select {}
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
