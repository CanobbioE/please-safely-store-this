// revive:disable:package-comments
package main

import (
	"fmt"
	"os"

	"github.com/CanobbioE/please-safely-store-this/cmd/psst"
)

func main() {
	if err := psst.RootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
