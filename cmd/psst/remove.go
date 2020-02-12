package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// removeCredentials deletes the file associated with the given account.
// The file removed is located at path/to/DefaultDir/.<account>
func removeCredentials(account string) {
	path := fmt.Sprintf("%s%s.%s", cfg.DefaultDir, string(filepath.Separator), account)
	err := os.Remove(filepath.FromSlash(path))
	if err != nil {
		log.Fatalf("Error removing credentials for %s: %v", account, err)
	}

	log.Infof("Removed %s", path)
}
