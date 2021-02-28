package encryption

import (
	"fmt"
	"log"
	"os"
	"io/ioutil"
	"crypto/aes"
    "crypto/cipher"
)

func Decrypt() {

	// file, err := os.Open(file)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Decrypts data
	key := []byte("passphrasewhichneedstobe32bytes!")
    ciphertext, err := ioutil.ReadFile("myfile.data")
    // if our program was unable to read the file
    // print out the reason why it can't
    if err != nil {
        log.Println(err)
    }

    // Deletes "myfile.data"
    er := os.Remove("myfile.data")
    if er != nil {
    	log.Fatal(er)
    }

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