package cryptoutils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"hash/fnv"
)

var iv = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

// GenerateKey32 hashes the given secret into a 32byte key.
// It uses fnv which is not a cryptographic hash function.
func GenerateKey32(secret string) string {
	hasher := fnv.New128a()
	hasher.Write([]byte(secret))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

// Encrypt takes a 32byte (128bit) key and a plain text.
// The key is used to encrypt the text using crypto/cipher
func Encrypt(key, text string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	plaintext := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	cfb.XORKeyStream(ciphertext, plaintext)
	return encodeBase64(ciphertext), nil
}

// Decrypt takes a 32byte (128bit) key and a plain text.
// The key is used to deencrypt the text using crypto/cipher
func Decrypt(key, text string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	ciphertext := decodeBase64(text)
	cfb := cipher.NewCFBEncrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	cfb.XORKeyStream(plaintext, ciphertext)
	return string(plaintext), nil
}

func encodeBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func decodeBase64(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}
