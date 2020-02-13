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
	Expecting(expectedInputs []string, msg string, onMismatch func()) PromptBuilder
}

type promptBuilder struct {
	basePrompter promptFunc
	modifier     modifyFunc
	in, out      *os.File
}

// ForSecret asks, on the specified channel,for a secret input.
// The in/out-put channels default to os.Stdin and os.Stdout.
// If you want to change these values call the On() function.
// The secret is not displayed while typed and it is not echoed
// after it is entered, on most terminals.
// On some terminals it fails (i.e. git bash), see https://github.com/golang/go/issues/34612
func ForSecret() PromptBuilder {
	fx := func(in, out *os.File, prompt string) (string, error) {
		oldState, _ := terminal.MakeRaw(int(in.Fd()))
		defer terminal.Restore(int(in.Fd()), oldState)

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

// ForConfirmation asks for confirmation
func ForConfirmation() PromptBuilder {
	fx := func(in, out *os.File, prompt string) (string, error) {
		fmt.Fprintf(out, prompt)
		rw := bufio.NewReadWriter(bufio.NewReader(in), bufio.NewWriter(out))
		t := terminal.NewTerminal(rw, "")
		return t.ReadLine()
	}
	return &promptBuilder{
		basePrompter: fx,
		in:           os.Stdin,
		out:          os.Stdout,
		modifier:     nil,
	}
}

// Expecting keep calling the base prompt function untill
// the user inputs an expected string.
// All the possible and expected inputs are passed through the expectedInputs parameter.
// If the input does not match any expected string, then onMismatch is called.
// The prompt msg ovewrite any message passed to the DoPrompt() funciton.
func (pb *promptBuilder) Expecting(expectedInputs []string, msg string, onMismatch func()) PromptBuilder {
	pb.modifier = func(fx promptFunc) (string, error) {
	repeat:
		in, err := fx(pb.in, pb.out, msg)
		if err != nil {
			return in, err
		}

		for _, s := range expectedInputs {
			if in == s {
				return s, nil
			}
		}
		onMismatch()
		goto repeat

	}
	return pb
}

// WithConfirm calls the base prompt function twice.
// If the two returend value are not equal onMismatch gets called and
// the function starts over.
// The two prompt messages override any message passed to the DoPrompt funciton.
// msgEnter is for the first prompt and msgConfirm for the second one.
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
