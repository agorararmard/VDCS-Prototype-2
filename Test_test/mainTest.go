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
	key1 := []byte("passphrasewhichneedstobe32bytes!")
	key2 := []byte("wrongphrasewhichneedstobe32bytes")
	c1, ok1 := EncryptAES(key1, text)
	//	c2 := EncryptAES(key2, text)
	if !ok1 {
		fmt.Println("I couldn't encrypt")
	}
	fmt.Println("Decryption Program v0.01")
	plaintext, ok2 := DecryptAES(key2, c1)

	if !ok2 {
		fmt.Println("I couldn't Decrypt")
	}
	plaintext, ok3 := DecryptAES(key1, c1)

	if !ok3 {
		fmt.Println("I couldn't Decrypt")
	}

	fmt.Println(string(plaintext))
}

func EncryptAES(encKey []byte, plainText []byte) (ciphertext []byte, ok bool) {

	ok = false //assume failure
	//			encKey = append(encKey, hash)
	c, err := aes.NewCipher(encKey)
	if err != nil {
		fmt.Println(err)
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		fmt.Println(err)
		return
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println(err)
		return
	}
	ciphertext = gcm.Seal(nonce, nonce, plainText, nil)
	//fmt.Println(ciphertext)
	ok = true

	return
}

func DecryptAES(encKey []byte, cipherText []byte) (plainText []byte, ok bool) {

	ok = false //assume failure

	c, err := aes.NewCipher(encKey)
	if err != nil {
		fmt.Println(err)
		return
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		fmt.Println(err)
		return
	}

	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		fmt.Println(err)
		return
	}

	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]
	plainText, err = gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(string(plaintext))
	ok = true
	return
}
