package main

import (
	"fmt"
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
	result string
	str string
	slice int
	iPlusShredSize int
	shredSize int = 4096
	counter = 0
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

	// Converts data into bytes
	b := []byte(readData)

    // Creates directory called "shredded_files"
    errDir := os.Mkdir("shredded_files", 0755)
    if errDir != nil {
            log.Fatal(errDir)
    }

    // Changes directory to "shredded_files"
    erro := os.Chdir("shredded_files")
    if erro != nil {
        log.Fatal(erro)
    }

	// Shreds file "myfile.data" and stores it in separate files
    for i := 0; i < len(b); i += shredSize {
    	iPlusShredSize = i + shredSize
    	str = strconv.Itoa(counter)
    	result = "file" + str + ".txt" // Names of shredded files
    	file, err := os.OpenFile(
        	result,
    	    os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
        	0666,
    	)
    	if err != nil {
        	log.Fatal(err)
    	}
    	defer file.Close()


    	if iPlusShredSize < len(b) {
    		bytesWritten, err := file.Write(b[i:iPlusShredSize])
    		if err != nil {
    			log.Fatal(err)
    		}
    		log.Printf("Wrote %d bytes.\n", bytesWritten)
    	} else {
    		slice = len(b) - i
    		bytesWritten, err := file.Write(b[i:i+slice])
    		if err != nil {
    			log.Fatal(err)
    		}
    		log.Printf("Wrote %d bytes.\n", bytesWritten)
    	}
    	counter++
	}	

    // Exits out of "shredded_files"
    errorr := os.Chdir("../")
    if errorr != nil {
        log.Fatal(errorr)
    }

	// Deletes "myfile.data"
	errr := os.Remove("myfile.data")
	if errr != nil {
		log.Fatal(errr)
	}
}







