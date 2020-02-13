package main

import (
	"bytes"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/atotto/clipboard"
)

// retrieveCredentials executes the logic for the psst -a <account> command.
// It decrypts the credentials file and prints the username to standard output.
// The password is copied in the system's clipboard, ready to be pasted anywhere.
func retrieveCredentials(account string) {
	path := fmt.Sprintf("%s%s.%s", cfg.DefaultDir, string(filepath.Separator), account)

	username, password := decryptCredentialsFile(path)

	// Copy password to clipboard
	if runtime.GOOS == "linux" {
		if !commandExists("xclip") && !commandExists("xsel") {
			log.Warnf("Copying to clipboard might fail, consider installing xclip or xsel.\n")
		}
	}
	b := bytes.Trim([]byte(password), "\000")
	if err := clipboard.WriteAll(string(b)); err != nil {
		log.Fatalf("Error copying password to clipboard: %v", err)
	}
	fmt.Printf("User: %s\n", username)

}
