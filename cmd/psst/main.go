package main

import (
	"flag"
	_log "log"
	"os"
	"path/filepath"

	"github.com/CanobbioE/please-safely-store-this/pkg/fileutils"
	"github.com/fatih/color"
)

const usage = `Usage:
	psst [-n | -r] -a ACCOUNT [-p PASSWORD, -u USERNAME]
	psst --config [-d DIRECTORY]

Options:
	-n, --new
		Used with -a/--account ACCOUNT and at least one between -p/--password PASSWORD
		and -u/--username USERNAME:
		creates a new credentials file at <default_directory>/.psst/.<ACCOUNT>,
		if the destination already exists it updates it with the given USERNAME and/or PASSWORD.

	-r, --remove
		Used with -a/--account ACCOUNT:
		remove a credentials file by deleting the file at <default_directory>/.psst/.<ACCOUNT>

	-a, --account ACCOUNT
		Used with -n/--new or -r/--r:
		specify the ACCOUNT onto which the operation is gonna take effect.
		Used by itself:
		retrive the ACCOUNT credentials.

	-p, --password PASSWORD
		Used with -n/--new:
		specify the path to the file from which the PASSWORD will be read
		and saved into the credentials file.
		The path can be both in Unix and in Windows format.

	-u, --username USERNAME
		Used with -n/--new:
		specify the USERNAME's value to be saved into the credentials file.

	-c, --config
		Used with -d/--directory:
		configure psst options.

	-d, --directory
		Used with -c/--config:
		specify the path to the DIRECTORY used to store the credentials file.
		The path could be both in Unix or in Windows format.
		The default path is <user_home>/.psst/

	-h, --help
		Show an helpful and well formatted message. :)

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

type PrintfFunc func(msg string, args ...interface{})
type Logger struct {
	Fatalf, Warnf, Infof PrintfFunc
}

var log = Logger{
	Fatalf: func(msg string, args ...interface{}) {
		_log.Printf(msg, args...)
		_log.Printf("If you feel like you need a refresh, try running psst --help to read the usage notes.")
		_log.Fatalf("If the error is something unexpected, please open an issue on CanobbioE/please-safely-store-this")
	},
	Infof: func(msg string, args ...interface{}) {
		_log.Printf(msg, args...)
	},
	Warnf: func(msg string, args ...interface{}) {
		color.Set(color.FgHiYellow)
		_log.Printf(msg, args...)
		color.Unset()
	},
}

func main() {
	// Set default dir
	savedDir := os.Getenv("PSSTDIR")

	if savedDir == "" {

		home, err := os.UserHomeDir()
		if err != nil {
			home = "./"
			log.Warnf("Error deriving user home directory, " +
				"creating one in the current directory\n")
		}
		savedDir = filepath.FromSlash(home + DEFAULTDIR)
	}

	onCreation := func() { log.Infof("Created default directory at %s\n", savedDir) }
	fileutils.CreateIfDoesntExist(savedDir, onCreation)

	cfg = ConfigOptions{
		DefaultDir: savedDir,
	}

	// Handle flags
	flag.Usage = func() { log.Infof("%s\n", usage) }
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
	flag.StringVar(&directoryFlag, "d", "", "default directory's value")
	flag.StringVar(&directoryFlag, "directory", "", "default directory's value")
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
