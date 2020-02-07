package fileutils

import "os"

// Exists returns true if the specified path exists
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// CreateIfDoesntExist checks for the path's existence,
// if it doesn't exists it creates all the directories that do not exists along
// that path and calls the onCreation function.
func CreateIfDoesntExist(path string, onCreation func()) error {
	ok, err := Exists(path)
	if err != nil {
		return err
	}
	if !ok {
		os.MkdirAll(path, os.ModePerm)
		onCreation()
	}
	return nil
}
