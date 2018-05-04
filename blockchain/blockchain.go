package blockchain

type BlockChain struct {
	blocks []Block
}

func NewBlockChain() BlockChain {
	return BlockChain{[]Block{NewGenesisBlock()}}
}

func (bc *BlockChain) AddBlock(data string) error {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(prevBlock, data)

	//if isBlockValid(newBlock, bc.blocks[len(bc.blocks)-1]) {
	newBlockchain := append(bc.blocks, newBlock)
	bc.blocks = newBlockchain

	return nil
}

func (bc *BlockChain) GetAllBlocks() []Block {
	return bc.blocks
}

// make sure using longest chain
func (bc *BlockChain) replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(bc.blocks) {
		bc.blocks = newBlocks
	}
}
