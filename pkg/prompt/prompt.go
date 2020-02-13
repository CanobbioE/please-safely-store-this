package prompt

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

// ForSecret asks, on the specified stream,
// for a secret input using the provided prompt.
// The secret is not displayed while typed and it is not echoed
// after it is entered, on most terminals.
// See https://github.com/golang/go/issues/34612
func ForSecret(in, out *os.File, prompt string) (string, error) {
	terminal.MakeRaw(int(in.Fd()))
	rw := bufio.NewReadWriter(bufio.NewReader(in), bufio.NewWriter(out))
	t := terminal.NewTerminal(rw, "")

	fmt.Fprintf(out, prompt)
	bytePassphrase, err := t.ReadPassword("")
	fmt.Fprintf(out, "\n")
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
func WithConfirm(in, out *os.File, promptEnter, promptConfirm string, onMismatch func()) (string, error) {
repeat:
	enter, err := ForSecret(in, out, promptEnter)
	if err != nil {
		return enter, err
	}

	confirm, err := ForSecret(in, out, promptConfirm)
	if err != nil {
		return enter, err
	}
	if string(enter) != string(confirm) {
		onMismatch()
		goto repeat
	}

	return string(enter), nil
}
