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
	Index        int
	Timestamp    string
	Transactions []*Transaction
	Hash         []byte
	PrevHash     []byte
	Nonce        int
}

// NewGenesisBlock generate GenesisBlock
func NewGenesisBlock(coinbase *Transaction) *Block {
	block := &Block{-1,
		"",
		nil,
		[]byte{},
		nil,
		0}
	return NewBlock(block, []*Transaction{coinbase})
}

// NewBlock generate new block using previous block's hash
func NewBlock(prevBlock *Block, transactions []*Transaction) *Block {

	newBlock := &Block{prevBlock.Index + 1,
		time.Now().String(),
		transactions,
		[]byte{},
		prevBlock.Hash,
		0,
	}

	pow := NewProofOfWork(newBlock)
	nonce, hash := pow.Run()

	newBlock.Hash = hash
	newBlock.Nonce = nonce

	return newBlock
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
	description := fmt.Sprintf("Prev hash: %x\nData: %s\nHash: %x\n", b.PrevHash, b.Transactions, b.Hash)
	return description
}
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}
