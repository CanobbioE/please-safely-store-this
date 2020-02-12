package cryptoutils

import (
	"crypto/aes"
	"crypto/cipher"
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
// The key is used to encrypt the text using authenticated encryption
func Encrypt(key, text string) (string, error) {
	plaintext := []byte(text)
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := []byte(key)[:aesgcm.NonceSize()]

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	return string(ciphertext), nil

}

// Decrypt takes a 32byte (128bit) key and a plain text.
// The key is used to decrypt the text using authenticated encryption
func Decrypt(key, text string) (string, error) {
	ciphertext := []byte(text)

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := []byte(key)[:aesgcm.NonceSize()]

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
