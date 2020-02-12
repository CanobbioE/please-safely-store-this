package cryptoutils

import (
	"testing"
)

func TestGenerateKey32(t *testing.T) {
	type args struct {
		secret string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"empty secret", args{""}, 32},
		{"single digit secret", args{"a"}, 32},
		{"multiple any char secret", args{"this is %v a secret!!123\000"}, 32},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateKey32(tt.args.secret); len(got) != tt.want {
				t.Errorf("GenerateKey32(%v) = len(%v) = %v, want %v", tt.args.secret, got, len(got), tt.want)
			}
		})
	}
}

func TestEncrypt(t *testing.T) {
	key32 := string([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32})
	type args struct {
		key  string
		text string
	}
	comparableOutput := func(a args) string {
		comp, err := Encrypt(a.key, a.text)
		if err != nil {
			panic(err)
		}
		return comp
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"should throw invalid key's length", args{"", "text"}, "", true},
		{"should throw invalid key's length", args{"short", "text"}, "", true},
		{"should encrypt empty text", args{key32, ""}, comparableOutput(args{key32, ""}), false},
		{"should encrypt any text", args{key32, "random %d text \x12!"}, comparableOutput(args{key32, "random %d text \x12!"}), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encrypt(tt.args.key, tt.args.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Encrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecrypt(t *testing.T) {
	key32 := string([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32})
	type args struct {
		key  string
		text string
	}
	encr := func(a args) string {
		r, err := Encrypt(a.key, a.text)
		if err != nil {
			panic(err)
		}
		return r
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"should throw invalid key's length", args{"", ""}, "", true},
		{"should throw failed authentication", args{key32, ""}, "", true},
		{"should throw invalid key's length", args{"short", ""}, "", true},
		{"should decrypt empty text", args{key32, encr(args{key32, ""})}, "", false},
		{"should decrypt any text", args{key32, encr(args{key32, "random %d text \x12!"})}, "random %d text \x12!", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decrypt(tt.args.key, tt.args.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Decrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}
