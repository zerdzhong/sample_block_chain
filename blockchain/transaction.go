package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

const subsidy = 20

// Transaction define contain id, Input, Output
type Transaction struct {
	ID     []byte
	Input  []TXInput
	Output []TXOutput
}

// NewCoinbaseTX create New coinbase transaction
func NewCoinbaseTX(to, data string) *Transaction {
	if "" == data {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	//coinbase 交易无需输入，矿工奖励，相当于货币发行
	trxIn := TXInput{[]byte{}, -1, data}
	trxOut := TXOutput{subsidy, to}

	trx := Transaction{nil, []TXInput{trxIn}, []TXOutput{trxOut}}
	trx.SetID()

	return &trx
}

// SetID sets ID for an transaction hash the transaction encode data
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

func (tx *Transaction) IsCoinBase() bool {
	return len(tx.Input) == 1 && len(tx.Input[0].TXID) == 0 && tx.Input[0].Vout == -1
}

// TXInput input of a transaction
type TXInput struct {
	TXID      []byte
	Vout      int
	ScriptSig string
}

// TXOutput output of a transaction
type TXOutput struct {
	Value        int
	ScriptPubKey string
}

// CanUnlockOutputWith checks whether the address initiated the transaction
func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}

// CanBeUnlockedWith checks if the output can be unlocked with the provided data
func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}

func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	acc, validOutputs := bc.FindSpendableOutputs(from, amount)

	if acc < amount {
		log.Panic("ERROR: 余额不足")
	}

	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)

		if nil != err {
			log.Panic(err)
		}

		for _, out := range outs {
			input := TXInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TXOutput{amount, to})
	if acc > amount {
		outputs = append(outputs, TXOutput{acc - amount, from}) // a change
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx
}
