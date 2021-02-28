package main

import (
	"fmt"
	"SecureBlock/blockchain"
	"os"
	"log"
	"io"
	"io/ioutil"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"strconv"
)

var (
	filename string
	shredSize int = 4096
	iPlusShredSize int
	slice int // extra bit of data left after chunking

)

func main() {
	fmt.Print("Please enter filename: ")
	fmt.Scanln(&filename)

	// Opens file inputted by user
	file, err := os.Open(filename)
    if err != nil { // handles errors
        log.Fatal(err)
    }

    // Reads data from the file
    data, err := ioutil.ReadAll(file)
    if err != nil {
        log.Fatal(err)
    }

    // Creates a key of size 32 bytes
    key := []byte("passphrasewhichneedstobe32bytes!")

    // generate a new aes cipher using our 32 byte long key
    c, err := aes.NewCipher(key)
    // if there are any errors, handle them
    if err != nil {
        fmt.Println(err)
    }

    // gcm or Galois/Counter Mode, is a mode of operation
    // for symmetric key cryptographic block ciphers
    // - https://en.wikipedia.org/wiki/Galois/Counter_Mode
    gcm, err := cipher.NewGCM(c)
    // if any error generating new GCM
    // handle them
    if err != nil {
        fmt.Println(err)
    }

    // creates a new byte array the size of the nonce
    // which must be passed to Seal
    nonce := make([]byte, gcm.NonceSize())
    // populates our nonce with a cryptographically secure
    // random sequence
    if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
        fmt.Println(err)
    }

    // here we encrypt our text using the Seal function
    // Seal encrypts and authenticates plaintext, authenticates the
    // additional data and appends the result to dst, returning the updated
    // slice. The nonce must be NonceSize() bytes long and unique for all
    // time, for a given key.
    // encryptedData := string())
    // the WriteFile method returns an error if unsuccessful
	err = ioutil.WriteFile("myfile.data", gcm.Seal(nonce, nonce, data, nil), 0777)
	// handle this error
	if err != nil {
  		// print it out
 	 	fmt.Println(err)
	}

	// Reads data from "myfile.data"
	readData, err := ioutil.ReadFile("myfile.data")
	if err != nil { // handles errors
		log.Fatal(err)
	}
	
	chain := blockchain.InitBlockChain()

	// Shreds file "myfile.data" and stores it in separate blocks
    for i := 0; i < len(readData); i += shredSize {
    	iPlusShredSize = i + shredSize
    	
    	if iPlusShredSize < len(readData) {
    		chain.AddBlock(readData[i:iPlusShredSize])
    	} else {
    		slice = len(readData) - i
    		chain.AddBlock(readData[i:i+slice])
    	}
	}

	for _, block := range chain.Blocks {
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Data in Block: %v\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}




}