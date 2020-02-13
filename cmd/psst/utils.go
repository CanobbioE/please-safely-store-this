package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/CanobbioE/please-safely-store-this/pkg/cryptoutils"
	"github.com/CanobbioE/please-safely-store-this/pkg/prompt"
)

// commandExists checks if a command exists or not
// Thanks to github.com/miguelmota.
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// readPlaintextPassword reads a password in plain text from a file
// and returns it.
func readPlaintextPassword(path string) string {
	password, err := ioutil.ReadFile(filepath.FromSlash(path))
	if err != nil {
		log.Fatalf("Error reading password file: %v", err)
	}
	return string(password)
}

// deleteFile deletes the specified file and handles any error.
func deletedFile(pathToPassword string) {
	err := os.Remove(pathToPassword)
	if err != nil {
		log.Fatalf("Error removing password's input file: %v", err)
	}
}

// decryptText returns the given text decrypted using the secret
// that the user provides from standard input.
func decryptText(text string) string {
	passphrase, err := prompt.ForSecret().On(os.Stdin, os.Stdout).DoPrompt("Encryption passphrase:")
	if err != nil {
		log.Fatalf("Error reading input: %v", err)
	}

	key := cryptoutils.GenerateKey32(passphrase)
	plaintext, err := cryptoutils.Decrypt(key, text)
	if err != nil {
		log.Fatalf("Error encrypting password: %v", err)
	}

	return plaintext
}

// encryptText returns the given text encrypted using the secret
// that the user provides from standard input.
func encryptText(text string) string {
	onMismatch := func() { log.Infof("Passphrases do not match.") }
	p := prompt.ForSecret().On(os.Stdin, os.Stdout)
	p = p.WithConfirm("Encryption passphrase:", "Confirm passphrase:", onMismatch)
	passphrase, err := p.DoPrompt("")
	if err != nil {
		log.Fatalf("Error reading input: %v", err)
	}

	key := cryptoutils.GenerateKey32(passphrase)
	encryptedText, err := cryptoutils.Encrypt(key, text)
	if err != nil {
		log.Fatalf("Error encrypting password: %v", err)
	}

	return encryptedText
}

// decryptCredentialsFile reads and returns the decrypted password and username from
// the specified file.
func decryptCredentialsFile(path string) (string, string) {
	fileContent, err := ioutil.ReadFile(filepath.FromSlash(path))
	if err != nil {
		log.Fatalf("Error reading from file: %v", err)
	}

	plaintext := decryptText(string(fileContent))

	lines := strings.Split(plaintext, "\n")
	return lines[0], lines[1]
}
