package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/CanobbioE/please-safely-store-this/pkg/fileutils"
)

// createOrUpdateCredentials calls createCredentials if there is no entry
// for the given account, otherwise it calls updateCredentials
func createOrUpdateCredentials(account, user, pathToPassword string) {
	pathToCredentials := fmt.Sprintf("%s%s.%s", DEFAULTDIR, string(filepath.Separator), account)
	updating, err := fileutils.Exists(pathToCredentials)
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
		deletedFile(pathToPassword)
	}

}

// createCredentials creates a new file at path/to/defaultDir/.<account>
// The generated file will contain the username in clear and the encrypted password.
func createCredentials(account, user, pathToPassword string) {
	if pathToPassword == "" || user == "" {
		log.Fatalf("Error: both -p/--password and -u/--user must be specified when creating with -n/--new.")
	}
	password := readPlaintextPassword(pathToPassword)
	_createOrUpdate(account, user, password)
}

// updateCredentials updates a file at path/to/defaultDir/.<account>
// If only one between password and user is specified,
// that's what is going to be updated.
func updateCredentials(account, user, pathToPassword string) {
	log.Infof("A credentials file for the specified account already exists, updating!")
	pathToCredentials := fmt.Sprintf("%s%s.%s", DEFAULTDIR, string(filepath.Separator), account)
	if pathToPassword == "" && user == "" {
		log.Fatalf("Error: at least one between -p/--password and -u/--user must be specified when updating with -n/--new.")
	}

	msg := fmt.Sprintf("Do you want to update %v? (y/n) ", pathToCredentials)
	if !userWantsToContinue(msg) {
		return
	}

	username, password := decryptCredentialsFile(pathToCredentials)
	log.Infof("Succesfully decrypted old file!")
	log.Infof("You can now change or keep the old passphrase.")

	switch {
	case pathToPassword != "":
		password = readPlaintextPassword(pathToPassword)
	default:
		user = username
	}

	_createOrUpdate(account, user, password)
}

// _createOrUpdate creates or update an account with the given credentials
// It encrypts the username and the password before saving them to file
func _createOrUpdate(account, user, password string) {
	pathToCredentials := filepath.FromSlash(fmt.Sprintf("%s/.%s", DEFAULTDIR, account))
	fileContent := fmt.Sprintf("%s\n%s", user, password)

	encryptedFileContent := encryptText(fileContent)

	err := ioutil.WriteFile(pathToCredentials, []byte(encryptedFileContent), os.ModePerm)
	if err != nil {
		log.Fatalf("Error creating credentials file: %v", err)
	}

	log.Infof("Added credential for user %s at %s\n", user, pathToCredentials)
}
