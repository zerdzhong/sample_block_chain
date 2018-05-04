package blockchain

import (
	"fmt"
	"github.com/boltdb/bolt"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"
const lastDBFileKey = "l"

// Blockchain DB key-value : "lastHash"

type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func NewBlockChain() *Blockchain {

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if nil == b {
			b, err = tx.CreateBucket([]byte(blocksBucket))
			genesisBlock := NewGenesisBlock()

			serializedData := genesisBlock.Serialize()

			err = b.Put(genesisBlock.Hash, serializedData)
			err = b.Put([]byte(lastDBFileKey), genesisBlock.Hash)

			tip = []byte(genesisBlock.Hash)
		} else {
			tip = b.Get([]byte(lastDBFileKey))
		}

		return nil
	})

	bc := Blockchain{tip, db}

	return &bc
}

func (bc *Blockchain) AddBlock(data string) error {

	var lastBlock *Block

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		encodedBlock := b.Get(bc.tip)
		lastBlock, _ = DeserializeBlock(encodedBlock)

		return nil
	})

	newBlock := NewBlock(lastBlock, data)

	bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		err = b.Put([]byte(lastDBFileKey), newBlock.Hash)

		bc.tip = newBlock.Hash

		return err
	})

	return err
}

func (bc *Blockchain) CloseDB() {
	bc.db.Close()
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}
	return bci
}

func (iterator *BlockchainIterator) Next() *Block {
	var block *Block

	err := iterator.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		encodeBlock := b.Get(iterator.currentHash)

		block, _ = DeserializeBlock(encodeBlock)

		return nil
	})

	if nil != err {
		fmt.Println(err.Error())
	}

	iterator.currentHash = block.PrevHash

	return block
}

//func (bc *Blockchain) GetAllBlocks() []Block {
//	//return bc.blocks
//}
