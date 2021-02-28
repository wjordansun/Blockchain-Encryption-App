package main

import (
	"fmt"
	"log"
	"os"
	"io/ioutil"
	"strconv"
	"crypto/aes"
    "crypto/cipher"
)

var (
	fileInfo os.FileInfo
	err error
	str string
	result string
	readData []byte
	bytesWritten int
	counter = 0
)

func main() {

	// Changes directory to "shredded_files"
	erro := os.Chdir("shredded_files")
    if erro != nil {
        log.Fatal(erro)
    }

    // Opens a file called "myfile.data"
	file, err := os.OpenFile(
		"myfile.data",
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
        0666,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Combines shredded files into one file and writes it to "myfile.data"
	for {
		str = strconv.Itoa(counter)
		result = "file" + str + ".txt"
		fileInfo, err = os.Stat(result)
		
		if err != nil {
			if os.IsNotExist(err) {
				break
			}
		}

		readData, err = ioutil.ReadFile(result)
		if err != nil { // handles errors
			log.Fatal(err)
		}

		bytesWritten, err = file.Write(readData)
    	if err != nil {
        	log.Fatal(err)
    	}
    	log.Printf("Wrote %d bytes.\n", bytesWritten)

   		err := os.Remove(result)
   		if err != nil {
   			log.Fatal(err)
   		}

		counter++
	}
	file.Close()

	// Decrypts data
	key := []byte("passphrasewhichneedstobe32bytes!")
    ciphertext, err := ioutil.ReadFile("myfile.data")
    // if our program was unable to read the file
    // print out the reason why it can't
    if err != nil {
        fmt.Println(err)
    }

    // Deletes "myfile.data"
    er := os.Remove("myfile.data")
    if er != nil {
    	log.Fatal(er)
    }

    // Exits out of "shredded_files"
    errorr := os.Chdir("../")
    if errorr != nil {
        log.Fatal(errorr)
    }

    // Removes "shredded_files"
    os.RemoveAll("shredded_files")

    // Continues decryption
    c, err := aes.NewCipher(key)
    if err != nil {
        fmt.Println(err)
    }

    gcm, err := cipher.NewGCM(c)
    if err != nil {
        fmt.Println(err)
    }

    nonceSize := gcm.NonceSize()
    if len(ciphertext) < nonceSize {
        fmt.Println(err)
    }

    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        fmt.Println(err)
    }

    // Writes decrypted text to "Reassembled.txt"
    errr := ioutil.WriteFile("Reassembled.txt", plaintext, 0666)
    if errr != nil {
        log.Fatal(errr)
    }

}





