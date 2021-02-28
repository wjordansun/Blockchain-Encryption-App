package main

import (
	"fmt"
	"FinalSecureBlock/blockchain"
	"FinalSecureBlock/encryption"
	"strconv"
	"os"
	"runtime"
	"flag"
	"log"
)

type CommandLine struct {
	blockchain *blockchain.BlockChain
}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" encrypt -file FILE_NAME - encrypt file")
	fmt.Println(" print - Prints the blocks in the chain")
	fmt.Println(" decrypt - decrypt file")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) encrypt(filename string) {
	const shredSize int = 4096
	var iMinusShredSize int

	data := encryption.Encrypt(filename)

	// Shreds file "myfile.data" and stores it in separate blocks
    for i := len(data); i > 0; i -= shredSize {
    	iMinusShredSize = i - shredSize
    	
    	if iMinusShredSize >= 0 {
    		cli.blockchain.AddBlock(data[iMinusShredSize:i])
    	} else {
    		cli.blockchain.AddBlock(data[:i])
    	}
	}

	errr := os.Remove("myfile.data")
	if errr != nil {
		log.Fatal(errr)
	}

	fmt.Println("File encrypted!")
}

func (cli *CommandLine) printChain() {
	iter := cli.blockchain.Iterator()

	for {
		block := iter.Next()

		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Data in Block: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break
		}
	}	
}

func (cli *CommandLine) decrypt() {
	iter := cli.blockchain.Iterator()

	file, err := os.OpenFile(
		"myfile.data",
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
        0666,
	)
	if err != nil {
		log.Fatal(err)
	}

	for {
		block := iter.Next()
		
		if string(block.Data) != "Genesis" {
			bytesWritten, err := file.Write(block.Data)
    		if err != nil {
       		 	log.Fatal(err)
    		}
    		log.Printf("Wrote %d bytes.\n", bytesWritten)
    	}

		if len(block.PrevHash) == 0 {

			break
		}
	}

	encryption.Decrypt()

}

func (cli *CommandLine) run() {
	cli.validateArgs()

	encryptCmd := flag.NewFlagSet("encrypt", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	decryptCmd := flag.NewFlagSet("decrypt", flag.ExitOnError)
	fileName := encryptCmd.String("file", " ", "fileName")

	switch os.Args[1] {
	case "encrypt":
		err := encryptCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "decrypt":
		err := decryptCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if encryptCmd.Parsed() {
		if *fileName == " " {
			encryptCmd.Usage()
			runtime.Goexit()
		}
		cli.encrypt(*fileName)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if decryptCmd.Parsed() {
		cli.decrypt()
	}
	
}

func main() {
	defer os.Exit(0)
	chain := blockchain.InitBlockChain()
	defer chain.Database.Close()

	cli := CommandLine{chain}
	cli.run()
}








