// source: https://gist.github.com/manishtpatel/8222606
package criptico

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// encrypt string to base64 crypto using AES
func Encrypt(key []byte, text string) string {
	// key := []byte(keyText)
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the cipher text.
	cipherText := make([]byte, aes.BlockSize+len(plaintext))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(cipherText)
}

// decrypt from base64 to decrypted string
func Decrypt(key []byte, cryptoText string) string {
	cipherText, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the cipher text.
	if len(cipherText) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(cipherText, cipherText)

	return fmt.Sprintf("%s", cipherText)
}