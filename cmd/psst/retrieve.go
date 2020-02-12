package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/CanobbioE/please-safely-store-this/pkg/cryptoutils"
	"github.com/CanobbioE/please-safely-store-this/pkg/prompt"
	"github.com/atotto/clipboard"
)

func retrieveCredentials(account string) {
	path := fmt.Sprintf("%s%s.%s", cfg.DefaultDir, string(filepath.Separator), account)
	username, encryptedPassword := readCredentialsFromFile(path)

	// Ask user for passphrase
	passphrase, err := prompt.ForSecret(os.Stdin, os.Stdout, "Encryption passphrase:")
	if err != nil {
		log.Fatalf("Error reading input: %v", err)
	}

	// Decrypt password
	key := cryptoutils.GenerateKey32(passphrase)
	password, err := cryptoutils.Decrypt(key, encryptedPassword)
	if err != nil {
		log.Fatalf("Error decrypting the password: %v", err)
	}

	// Copy password to clipboard
	if runtime.GOOS == "linux" {
		if !commandExists("xclip") && !commandExists("xsel") {
			log.Warnf("Copying to clipboard might fail, consider installing xclip or xsel.\n")
		}
	}
	b := bytes.Trim([]byte(password), "\000")
	log.Infof(string(b))
	if err := clipboard.WriteAll(password); err != nil {
		log.Fatalf("Error copying password to clipboard: %v", err)
	}
	fmt.Printf("User: %s\n", username)

}

// commandExists checks if a command exists or not
// Thanks to github.com/miguelmota
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// readCredentialsFromFile reads and returns an encrypted password and a username from
// the specified file.
func readCredentialsFromFile(path string) (user, encryptedPassword string) {
	fileContent, err := ioutil.ReadFile(filepath.FromSlash(path))
	if err != nil {
		log.Fatalf("Error reading from file: %v", err)
	}

	// since fileContent is just two lines, we can use Split
	lines := strings.Split(string(fileContent), "\n")
	user, encryptedPassword = lines[0], lines[1]
	return
}
