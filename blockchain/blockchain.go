package blockchain

import "errors"

type BlockChain struct {
	blocks []Block
}

func NewBlockChain() BlockChain {
	return BlockChain{[]Block{NewGenesisBlock()}}
}

func (bc *BlockChain) AddBlock(data string) error {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(prevBlock, data)

	if isBlockValid(newBlock, bc.blocks[len(bc.blocks)-1]) {
		newBlockchain := append(bc.blocks, newBlock)
		bc.replaceChain(newBlockchain)

		return nil
	}

	return errors.New("Block not valid")
}

// make sure using longest chain
func (bc *BlockChain) replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(bc.blocks) {
		bc.blocks = newBlocks
	}
}
