package blockchain

// Take the data from the block

// create a counter (nonce) which starts at 0

// create a hash of the data plus the counter

// check the hash to see if it meets a set of requirements

// Requirements:
// The First few bytes must contain 0's

import (
	"fmt"
	"log"
	"bytes"
	"crypto/sha256"
	"math"
	"math/big"
	"encoding/binary"
)

const Difficulty = 18 // Difficulty can increase based on amount of miners

type ProofOfWork struct {
	Block *Block
	Target *big.Int // requirement
}

func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty)) // Left shift // 256 is # of bytes in hash

	pow := &ProofOfWork{b, target}

	return pow
}

func (pow *ProofOfWork) InitData(nonce int) []byte {
	data := bytes.Join( // Creates a cohesive set of bytes
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.Data,
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)
	return data
}

func (pow *ProofOfWork) Run() (int, []byte) { // returns nonce and hash
	var intHash big.Int // type of int that can store up to 8 bytes
	var hash [32]byte

	nonce := 0

	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
		hash = sha256.Sum256(data)

		fmt.Printf("\r%x", hash)
		intHash.SetBytes(hash[:]) // converts hash to a bigInt to compare with target

		if intHash.Cmp(pow.Target) == -1 { // hash is less than target which means
			break							// block has been signed
		} else {
			nonce++
		}
	}
	fmt.Println()

	return nonce, hash[:]

}

func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int

	data := pow.InitData(pow.Block.Nonce)

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}

func ToHex(num int64) [] byte { // changes int to bytes
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)// BigEndian says how to organize 
		if err != nil {								// the bytes
			log.Panic(err)
		}

	return buff.Bytes()

}











