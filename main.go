package main

import (
	"fmt"
	"github.com/CanobbioE/please-safely-store-this/cmd/psst"
	"os"
)

func main() {
	if err := psst.RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
