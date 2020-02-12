package fileutils

import (
	"testing"
)

func TestExists(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Exists(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Exists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateIfDoesntExist(t *testing.T) {
	type args struct {
		path       string
		onCreation func()
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateIfDoesntExist(tt.args.path, tt.args.onCreation); (err != nil) != tt.wantErr {
				t.Errorf("CreateIfDoesntExist() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
