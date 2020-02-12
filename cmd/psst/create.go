package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/CanobbioE/please-safely-store-this/pkg/cryptoutils"
	"github.com/CanobbioE/please-safely-store-this/pkg/fileutils"
	"github.com/CanobbioE/please-safely-store-this/pkg/prompt"
)

// createOrUpdateCredentials calls createCredentials if there is no entry
// for the given account, otherwise it calls updateCredentials
func createOrUpdateCredentials(account, user, pathToPassword string) {
	updating, err := fileutils.Exists(cfg.DefaultDir + "/.account")
	if err != nil {
		log.Fatalf("Error checking for credentials existence: %v", err)
	}
	switch {
	case updating:
		updateCredentials(account, user, pathToPassword)
	default: // creating
		createCredentials(account, user, pathToPassword)
	}
	if pathToPassword != "" {
		deleteInputPasswordFile(pathToPassword)
	}

}

// createCredentials creates a new file at path/to/defaultDir/.<account>
// The generated file will contain the username in clear and the encrypted password.
func createCredentials(account, user, pathToPassword string) {
	if pathToPassword == "" || user == "" {
		log.Fatalf("Error: both -p/--password and -u/--user must be specified when creating with -n/--new.")
	}
	encryptedPassword := readAndEncryptPassword(pathToPassword)
	_createOrUpdate(account, user, encryptedPassword)
}

// updateCredentials updates a file at path/to/defaultDir/.<account>
// If only one between password and user is specified,
// that's what is going to be updated.
func updateCredentials(account, user, pathToPassword string) {
	log.Infof("A file for the specified account already exists, updating!")
	if pathToPassword == "" && user == "" {
		log.Fatalf("Error: at least one between -p/--password and -u/--user must be specified when updating with -n/--new.")
	}

	path := fmt.Sprintf("%s%s.%s", cfg.DefaultDir, string(filepath.Separator), account)
	username, encryptedPassword := readCredentialsFromFile(path)

	switch {
	case pathToPassword != "":
		encryptedPassword = readAndEncryptPassword(pathToPassword)
	default:
		user = username
	}

	_createOrUpdate(account, user, encryptedPassword)
}

// encryptPassword returns the password encrypted using the secret read from
// standard input
func encryptPassword(password string) string {
	onMismatch := func() { log.Infof("Passphrases do not match.") }
	passphrase, err := prompt.WithConfirm("Encryption passphrase:", "Confirm passphrase:", onMismatch)
	if err != nil {
		log.Fatalf("Error reading input: %v", err)
	}
	key := cryptoutils.GenerateKey32(passphrase)
	encryptedPassword, err := cryptoutils.Encrypt(key, string(password))
	if err != nil {
		log.Fatalf("Error encrypting password: %v", err)
	}
	return encryptedPassword
}

// writeCredentialsToFile creates or updates a file with the specified credentials
func writeCredentialsToFile(user, password, path string) {
	fileContent := fmt.Sprintf("%s\n%s", user, password)
	pathToCredentials := filepath.FromSlash(path)
	err := ioutil.WriteFile(pathToCredentials, []byte(fileContent), os.ModePerm)
	if err != nil {
		log.Fatalf("Error creating credentials file: %v", err)
	}
}

// readAndEncryptPassword reads a password in plain text from a file and
// returns the encrypted password, after deleting the input file.
func readAndEncryptPassword(path string) string {
	password, err := ioutil.ReadFile(filepath.FromSlash(path))
	if err != nil {
		log.Fatalf("Error reading password file: %v", err)
	}
	return encryptPassword(string(password))
}

// deleteInputPasswordFile deletes the file containing the clear password
func deleteInputPasswordFile(pathToPassword string) {
	err := os.Remove(pathToPassword)
	if err != nil {
		log.Fatalf("Error removing password's input file: %v", err)
	}
}

// _createOrUpdate creates or update an account with the given credentials
func _createOrUpdate(account, user, password string) {
	path := fmt.Sprintf("%s/.%s", cfg.DefaultDir, account)
	writeCredentialsToFile(user, password, path)
	log.Infof("Added credential for user %s at %s\n", user, path)
}
