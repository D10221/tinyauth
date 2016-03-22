// source: https://gist.github.com/manishtpatel/8222606
package tinyauth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"bytes"
	"errors"
)

type Criptico interface {
	Encrypt(text string) (string, error)
	Decrypt(cryptoText string) (string , error )
}

type DefaultCriptico struct {
	Key string
}

func (c *DefaultCriptico) Encrypt(text string ) (string,error) {
	b:=bytes.NewBufferString(c.Key).Bytes()
	return encrypt(b, text)
}

func (c *DefaultCriptico) Decrypt(text string ) (string, error) {
	return decrypt(bytes.NewBufferString(c.Key).Bytes(), text)
}

// encrypt string to base64 crypto using AES
func encrypt(key []byte, text string) (string, error) {
	// key := []byte(keyText)
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the cipher text.
	cipherText := make([]byte, aes.BlockSize+len(plaintext))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(cipherText) , nil
}

// decrypt from base64 to decrypted string
func decrypt(key []byte, cryptoText string) (string, error) {
	cipherText, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the cipher text.
	if len(cipherText) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(cipherText, cipherText)

	return fmt.Sprintf("%s", cipherText), nil
}