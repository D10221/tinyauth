// source: https://gist.github.com/manishtpatel/8222606
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"errors"
	"log"
)

type Criptico interface {
	Encrypt(text string) (string, error)
	Decrypt(cryptoText string) (string , error )
}

type DefaultCriptico struct {
	Key string
}

// encrypt string to base64 crypto using AES
func (c *DefaultCriptico) Encrypt(text string ) (string,error) {

	// key :=bytes.NewBufferString(c.Key).Bytes()
	key := []byte(c.Key)
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
func (c *DefaultCriptico) Decrypt(cryptoText string ) (string, error) {

	//key :=bytes.NewBufferString(c.Key).Bytes()
	key := []byte(c.Key)
	cipherText, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the cipher text.
	if len(cipherText) < aes.BlockSize {
		log.Printf("cipherText: %s ", cipherText)
		return "", errors.New("ciphertext too short")
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(cipherText, cipherText)

	return fmt.Sprintf("%s", cipherText), nil
}

