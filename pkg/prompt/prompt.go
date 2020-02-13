package prompt

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

type promptFunc func(in, out *os.File, msg string) (string, error)
type modifyFunc func(fx promptFunc) (string, error)

type PromptBuilder interface {
	On(in, out *os.File) PromptBuilder
	WithConfirm(promptEnter, promptConfirm string, onMismatch func()) PromptBuilder
	DoPrompt(msg string) (string, error)
}

type promptBuilder struct {
	basePrompter promptFunc
	modifier     modifyFunc
	in, out      *os.File
}

// ForSecret asks, on the specified stream,
// for a secret input using the provided prompt.
// The secret is not displayed while typed and it is not echoed
// after it is entered, on most terminals.
// The in/out-put channels default to os.Stdin and os.Stdout.
// See https://github.com/golang/go/issues/34612
func ForSecret() PromptBuilder {
	fx := func(in, out *os.File, prompt string) (string, error) {
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

	return &promptBuilder{
		basePrompter: fx,
		in:           os.Stdin,
		out:          os.Stdout,
		modifier:     nil,
	}
}

// WithConfirm keeps asking on the specified stream,
// for a secret input and a confirmation input
// untill the two inputs are equal.
// If the two inputs do not match the onMismatch function is called.
// The secrets are not displayed while typed and they are not echoed
// after they are entered, on most terminals.
// The two prompt messages override any message passed to the DoPrompt funciton.
// See https://github.com/golang/go/issues/34612
func (pb *promptBuilder) WithConfirm(msgEnter, msgConfirm string, onMismatch func()) PromptBuilder {
	pb.modifier = func(fx promptFunc) (string, error) {
	repeat:
		enter, err := fx(pb.in, pb.out, msgEnter)
		if err != nil {
			return enter, err
		}

		confirm, err := fx(pb.in, pb.out, msgConfirm)
		if err != nil {
			return enter, err
		}
		if string(enter) != string(confirm) {
			onMismatch()
			goto repeat
		}

		return string(enter), nil
	}

	return pb
}

// DoPrompt execute the full prompt sequence applying any modifier if present.
// The msg is used only with the base prompter, which means it is
// overwritten by any modifier that accepts a prompt message
func (pb *promptBuilder) DoPrompt(msg string) (string, error) {
	if pb.modifier != nil {
		return pb.modifier(pb.basePrompter)
	}
	return pb.basePrompter(pb.in, pb.out, msg)
}

// On sets the in/out-put channel for the prompt.
func (pb *promptBuilder) On(in, out *os.File) PromptBuilder {
	pb.in = in
	pb.out = out
	return pb
}
