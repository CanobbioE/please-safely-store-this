package fileutils

import (
	"os"
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
		{"should return true", args{"."}, true, false},
		{"should return false", args{"./thisfiledoesnotexist"}, false, false},
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
	emptyFunc := func() { return }

	didWorkCreation := func() (bool, string) {
		if ok, _ := Exists("./deleteme"); !ok {
			return false, "Did not create file"
		}
		return true, ""
	}

	defer os.RemoveAll("./deleteme")
	type args struct {
		path       string
		onCreation func()
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		didWork func() (bool, string)
	}{
		{"should create a file", args{"./deleteme", emptyFunc}, false, didWorkCreation},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateIfDoesntExist(tt.args.path, tt.args.onCreation); (err != nil) != tt.wantErr {
				t.Errorf("CreateIfDoesntExist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if did, msg := tt.didWork(); !did {
				t.Errorf(msg)
			}
		})
	}
}
