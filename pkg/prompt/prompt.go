package prompt

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

// ForSecret asks, on the specified stream,
// for a secret input using the provided prompt.
// The secret is not displayed while typed and it is not echoed
// after it is entered, on most terminals.
// See https://github.com/golang/go/issues/34612
func ForSecret(stream *os.File, prompt string) (string, error) {
	fmt.Print(prompt)
	terminal.MakeRaw(int(stream.Fd()))
	t := terminal.NewTerminal(io.ReadWriter(stream), "")
	bytePassphrase, err := t.ReadPassword("")
	fmt.Println("")
	if err != nil {
		return string(bytePassphrase), err
	}
	return string(bytePassphrase), nil
}

// WithConfirm keeps asking on the specified stream,
// for a secret input and a confirmation input
// untill the two inputs are equal.
// If the two inputs do not match the onMismatch function is called.
// The secrets are not displayed while typed and they are not echoed
// after they are entered, on most terminals.
// See https://github.com/golang/go/issues/34612
func WithConfirm(stream *os.File, promptEnter, promptConfirm string, onMismatch func()) (string, error) {
repeat:
	enter, err := ForSecret(stream, promptEnter)
	if err != nil {
		return enter, err
	}

	confirm, err := ForSecret(stream, promptConfirm)
	if err != nil {
		return enter, err
	}

	if string(enter) != string(confirm) {
		onMismatch()
		goto repeat
	}
	return string(enter), nil
}
