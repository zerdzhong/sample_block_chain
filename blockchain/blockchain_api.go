package blockchain

import (
	"fmt"
	"strconv"
)

var SampleChain = NewBlockchain()

func GetBalance(address string) int {

	balance := 0
	UTXOs := SampleChain.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)

	return balance
}

func Send(from, to string, amount int) {

	if 0 == amount {
		println("无效值")
		return
	}

	tx := NewUTXOTransaction(from, to, amount, SampleChain)
	SampleChain.MineBlock([]*Transaction{tx})
	fmt.Println("Success!")
}

func GetBlockChains() []*Block {
	iterator := SampleChain.Iterator()
	var chain []*Block

	for {
		block := iterator.Next()

		chain = append(chain, block)
		if len(block.PrevHash) == 0 {
			break
		}
	}

	return chain
}

func ReplaceBlockChain(newChain []Block) {
	chain := GetBlockChains()
	if len(newChain) > len(chain) {
		SampleChain.ReplaceChain(newChain)
	}
}

func PrintChain() {

	bci := SampleChain.Iterator()

	for {
		block := bci.Next()

		fmt.Printf(block.Description())
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break
		}
	}

}
