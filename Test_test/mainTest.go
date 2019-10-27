package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

func main() {
	fmt.Println("Encryption Program v0.01")

	text := []byte("My Super Secret Code Stuff")
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

	ciphertext := gcm.Seal(nonce, nonce, text, nil)
	fmt.Println(ciphertext)
	fmt.Println("Decryption Program v0.01")

	//  key := []byte("passphrasewhichneedstobe32bytes!")
	//, err := ioutil.ReadFile("myfile.data")
	// if our program was unable to read the file
	// print out the reason why it can't
	//	if err != nil {
	//		fmt.Println(err)
	//	}

	ce, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
	}

	gcme, err := cipher.NewGCM(ce)
	if err != nil {
		fmt.Println(err)
	}

	nonceSize := gcme.NonceSize()
	if len(ciphertext) < nonceSize {
		fmt.Println(err)
	}

	noncea, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcme.Open(nil, noncea, ciphertext, nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(plaintext))
}
