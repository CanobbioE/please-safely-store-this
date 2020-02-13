package main

import (
	"os"
	"path/filepath"

	"github.com/CanobbioE/please-safely-store-this/pkg/fileutils"
	_dir "github.com/otiai10/copy"
)

// ConfigOptions represents the configuration options for psst
type ConfigOptions struct {
	DefaultDir string
}

// configurePsst changes the default directory into which credential files are stored.
// Everything in the previous folder gets copied to the new one, and the previous one is deleted.
func configurePsst(path string) {
	// X-platform path
	path = filepath.FromSlash(path)

	// Create the folder if it doesn't exist
	onCreation := func() { log.Infof("Created directory at %s\n", path) }
	if err := fileutils.CreateIfDoesntExist(path, onCreation); err != nil {
		log.Fatalf("Error while creating the direcory %s: %v", path, err)
	}

	// Move everything from the old direcrory to the new one
	if err := _dir.Copy(cfg.DefaultDir, path); err != nil {
		log.Fatalf("Error while copying files from %s to %s : %v", err, cfg.DefaultDir, path)
	}
	log.Infof("Copied all files from %s to %s\n", cfg.DefaultDir, path)

	// Delete the old folder
	if err := os.Remove(cfg.DefaultDir); err != nil {
		log.Warnf("Couldn't remove old folder %s: %v", cfg.DefaultDir, err)
	}

	// Change settings
	os.Setenv("PSSTDIR", path)
	cfg.DefaultDir = path
}
