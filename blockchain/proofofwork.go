package blockchain

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"samplechain/utils"
)

//
const targetBits = 8

var (
	maxNonce = math.MaxInt64
)

// ProofOfWork struct define
type ProofOfWork struct {
	block  Block
	target *big.Int
}

// NewProofOfWork create new proof of work
func NewProofOfWork(b Block) ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := ProofOfWork{b, target}

	return pow
}

// Run mining hash
func (pow ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte

	nonce := 0

	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Transactions)

	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}

	fmt.Println("\nMining done")

	return nonce, hash[:]
}

// Validate is a block by pow
func (pow ProofOfWork) Validate() bool {
	var hashInt big.Int
	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}

func (pow ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			[]byte(pow.block.PrevHash),
			[]byte(pow.block.HashTransactions()),
			[]byte(pow.block.Timestamp),
			utils.IntToHex(targetBits),
			utils.IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}
