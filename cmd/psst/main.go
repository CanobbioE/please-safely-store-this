package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/CanobbioE/please-safely-strore-this/pkg/fileutils"
)

const usage = `Usage:
	psst [-n | -r] -a ACCOUNT [-p PASSWORD, -u USERNAME]
	psst --config [-d DIRECTORY]

Options:
	-n, --new                       Store or update a set of credentials, expects to be used with --account and --username and/or --password.
	-r, --remove                    Remove a credentials set, expects to be used with --account.
	-a, --account ACCOUNT           Specify the ACCOUNT's value, to be used with --new or --remove or alone.
	-p, --password PASSWORD         Specify the path to the file containing the password in clear, to be used with --new.
	-u, --username USERNAME         Specify the USERNAME's value, to be used with --new.
	-c, --config                    Configure the command.
	-d, --directory                 Specify the path to the DIRECTORY into which store/read the encrypted credentials, to be used with --config.
	-h, --help                      Show this message.

DIRECTORY defaults to ~/.psst in no other definition is found in the current environment.

Calls with no OPTIONS retrieve the credentials for the given ACCOUNT.


Example:
	$ psst -n -a grandma_instagram -p ./password.txt -u example@example.com
	Encryption passphrase:
	Confirm passphrase:
	Added credential for user example@example.com at ~/.psst/.grandma_instagram

	$ psst --config -d ~/myFolder
	Moved all the credentials from ~/.psst to ~/myFolder

	$ psst -a grandma_instagram
	Encryption passphrase:
	User: example@example.com
	Password copied to clipboard

	$ psst -r -a grandma_instagram
	Removed ~/myFolder/.grandma_instagram
`

const DEFAULTDIR = "/.psst"

var cfg ConfigOptions

func main() {
	// Set default dir
	savedDir := os.Getenv("PSSTDIR")

	if savedDir == "" {

		home, err := os.UserHomeDir()
		if err != nil {
			home = "./"
			log.Println("Error deriving user home directory, " +
				"creating one in the current directory")
		}
		savedDir = filepath.FromSlash(home + DEFAULTDIR)
	}

	onCreation := func() { log.Printf("Created default directory at %s", savedDir) }
	fileutils.CreateIfDoesntExist(savedDir, onCreation)

	cfg = ConfigOptions{
		DefaultDir: savedDir,
	}

	// Handle flags
	flag.Usage = func() { fmt.Fprintf(os.Stderr, "%s\n", usage) }
	var (
		configFlag, newFlag, removeFlag                        bool
		accountFlag, usernameFlag, passwordFlag, directoryFlag string
	)

	flag.BoolVar(&configFlag, "c", false, "configure the command")
	flag.BoolVar(&configFlag, "config", false, "configure the command")
	flag.BoolVar(&newFlag, "n", false, "create or update credentials")
	flag.BoolVar(&newFlag, "new", false, "create or update a set of credentials")
	flag.BoolVar(&removeFlag, "r", false, "remove a set of credentials")
	flag.BoolVar(&removeFlag, "remove", false, "remove a set of credentials")
	flag.StringVar(&usernameFlag, "u", "", "username's value")
	flag.StringVar(&usernameFlag, "username", "", "username's value")
	flag.StringVar(&passwordFlag, "p", "", "password's value")
	flag.StringVar(&passwordFlag, "password", "", "password's value")
	flag.StringVar(&accountFlag, "a", "", "account's name")
	flag.StringVar(&accountFlag, "account", "", "account's name")
	flag.StringVar(&directoryFlag, "d", cfg.DefaultDir, "default directory's value")
	flag.StringVar(&directoryFlag, "directory", cfg.DefaultDir, "default directory's value")
	flag.Parse()

	// Check flags correctness
	switch {
	case configFlag:
		switch {
		case newFlag:
			log.Fatalf("Error: -n/--new cannot be used with -c/--config.")
		case removeFlag:
			log.Fatalf("Error: -r/--remove cannot be used with -c/--config.")
		case accountFlag != "":
			log.Fatalf("Error: -a/--account cannot be used with -c/--config.")
		case usernameFlag != "":
			log.Fatalf("Error: -u/--user cannot be used with -c/--config.")
		case passwordFlag != "":
			log.Fatalf("Error: -p/--password cannot be used with -c/--config.")
		case directoryFlag == "":
			log.Fatalf("Error: -d/--directory must be specified with -c/--config.")
		}
		configurePsst(directoryFlag)
	case newFlag:
		if accountFlag == "" {
			log.Fatalf("Error: -a/--account must be specified with -n/--new.")
		}
		createOrUpdateCredentials(accountFlag, usernameFlag, passwordFlag)
	case removeFlag:
		if accountFlag == "" {
			log.Fatalf("Error: -a/--account must be specified with -r/--remove.")
		}
		removeCredentials(accountFlag)
	default: // retrieve credentials
		if accountFlag == "" {
			log.Fatalf("Error: -a/--account must be specified.")
		}
		retrieveCredentials(accountFlag)
	}
}
