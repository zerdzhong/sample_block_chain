package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"time"
)

// Block represent each item in blockchain
type Block struct {
	Index     int
	Timestamp string
	Data      string
	Hash      []byte
	PrevHash  []byte
	Nonce     int
}

// NewGenesisBlock generate GenesisBlock
func NewGenesisBlock() Block {
	t := time.Now()
	genesisBlock := Block{0, t.String(), "Genesis Block", []byte{}, []byte{}, 0}
	genesisBlock.Hash = calculateHash(genesisBlock)
	return genesisBlock
}

// NewBlock generate new block using previous block's hash
func NewBlock(prevBlock *Block, data string) *Block {
	var newBlock Block

	newBlock.Data = data

	newBlock.Index = prevBlock.Index + 1
	newBlock.Timestamp = time.Now().String()
	newBlock.PrevHash = prevBlock.Hash

	pow := NewProofOfWork(newBlock)
	nonce, hash := pow.Run()

	newBlock.Hash = hash
	newBlock.Nonce = nonce

	return &newBlock
}

//Serialize to byte save in DB
func (b *Block) Serialize() []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	encoder.Encode(b)

	return result.Bytes()
}

//DeserializeBlock from byte in DB
func DeserializeBlock(d []byte) (*Block, error) {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)

	return &block, err
}

//Description string description of block
func (b *Block) Description() string {
	description := fmt.Sprintf("Prev hash: %x\nData: %s\nHash: %x\n", b.PrevHash, b.Data, b.Hash)
	return description
}

// calculateHas SHA256 hashing
func calculateHash(block Block) []byte {
	blockInfo := string(block.Index) + block.Timestamp + block.Data + string(block.PrevHash)

	sha256 := sha256.New()
	sha256.Write([]byte(blockInfo))
	hashed := sha256.Sum(nil)

	return hashed
}
