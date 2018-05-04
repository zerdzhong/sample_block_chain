package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// Block represent each item in blockchain
type Block struct {
	Index     int
	Timestamp string
	Data      string
	Hash      string
	PrevHash  string
}

func NewGenesisBlock() Block {
	t := time.Now()
	genesisBlock := Block{0, t.String(), "Genesis Block", "", ""}
	return genesisBlock
}

// generate new block using previous block's hash
func NewBlock(prevBlock Block, data string) Block {
	var newBlock Block

	newBlock.Data = data

	newBlock.Index = prevBlock.Index + 1
	newBlock.Timestamp = time.Now().String()
	newBlock.PrevHash = prevBlock.Hash
	newBlock.Hash = calculateHash(newBlock)

	return newBlock
}

// calculateHas SHA256 hashing
func calculateHash(block Block) string {
	blockInfo := string(block.Index) + block.Timestamp + block.Data + block.PrevHash

	sha256 := sha256.New()
	sha256.Write([]byte(blockInfo))
	hashed := sha256.Sum(nil)

	return hex.EncodeToString(hashed)
}

// check block is valid
func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}
	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}
	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}
	return true
}
