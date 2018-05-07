package blockchain

import (
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"
const lastDBFileKey = "l"

const genesisCoinbaseData = "Genesis coinbase data"

// Blockchain DB key-value : "lastHash"

// Blockchain contain db and tip
type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

// Iterator iterator of blockchain
type Iterator struct {
	currentHash []byte
	db          *bolt.DB
}

// NewBlockchain create Blockchain
func NewBlockchain(address string) *Blockchain {

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if nil == b {

			genesisTransaction := NewCoinbaseTX(address, genesisCoinbaseData)
			genesisBlock := NewGenesisBlock(genesisTransaction)
			serializedData := genesisBlock.Serialize()

			b, err = tx.CreateBucket([]byte(blocksBucket))
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

// MineBlock add a block to Blockchain
func (bc *Blockchain) MineBlock(transactions []*Transaction) error {

	var lastBlock *Block

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		encodedBlock := b.Get(bc.tip)
		lastBlock, _ = DeserializeBlock(encodedBlock)

		return nil
	})

	newBlock := NewBlock(lastBlock, transactions)

	bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		err = b.Put([]byte(lastDBFileKey), newBlock.Hash)

		bc.tip = newBlock.Hash

		return err
	})

	return err
}

//FindUnSpentTransactions UTXO 寻找特定地址未被话费的输出（结余？）
func (bc *Blockchain) FindUnSpentTransactions(address string) []Transaction {
	var unSpentTXs []Transaction
	spentTXOutputs := make(map[string][]int)

	bci := bc.Iterator()

	//遍历所有区块
	for {
		block := bci.Next()
		for _, trx := range block.Transactions {
			trxID := hex.EncodeToString(trx.ID)

		FindOutPuts:
			for outIdx, out := range trx.Output {
				//outputs 已经花掉
				if nil != spentTXOutputs[trxID] {
					for _, spentOutput := range spentTXOutputs[trxID] {
						if spentOutput == outIdx {
							continue FindOutPuts
						}
					}
				}

				//找到转给 address 的 outputs
				if out.CanBeUnlockedWith(address) {
					unSpentTXs = append(unSpentTXs, *trx)
				}
			}

			if false == trx.IsCoinBase() {
				for _, input := range trx.Input {
					if !input.CanUnlockOutputWith(address) {
						continue
					}
					inputTRXID := hex.EncodeToString(input.TXID)
					spentTXOutputs[inputTRXID] = append(spentTXOutputs[inputTRXID], input.Vout)
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return unSpentTXs
}

func (bc *Blockchain) FindUTXO(address string) []TXOutput {
	var UTXOs []TXOutput

	unspentTransaction := bc.FindUnSpentTransactions(address)

	for _, tx := range unspentTransaction {
		for _, out := range tx.Output {
			if out.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

func (bc *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTXs := bc.FindUnSpentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)
		for outIdx, out := range tx.Output {
			if out.CanBeUnlockedWith(address) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
}

// CloseDB close DB
func (bc *Blockchain) CloseDB() {
	bc.db.Close()
}

// Iterator get blockchain's iterator
func (bc *Blockchain) Iterator() *Iterator {
	bci := &Iterator{bc.tip, bc.db}
	return bci
}

//Next of iterator
func (iterator *Iterator) Next() *Block {
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
