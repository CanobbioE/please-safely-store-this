# Please Safely Store This (PSST)
Psst is a password manager. There are many password managers out there but I wanted to create my own.

(Yes I did it only because I came up with the fun acronym)

## Installation
Binaries are provided [here](https://github.com/CanobbioE/please-safely-store-this/releases).

Or you can use Go CLI
```bash
$ go get github.com/CanobbioE/please-safely-store-this

$ go install ./CanobbioE/please-safely-store-this/cmd/
```

## Usage
```
Usage:
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

	-h, --help
		Show an helpful and well formatted message. :)

Example:
	$ psst -n -a grandma_instagram -p ./password.txt -u example@example.com
	Encryption passphrase:
	Confirm passphrase:
	Added credential for user example@example.com at ~/.psst/.grandma_instagram

	$ psst -a grandma_instagram
	Encryption passphrase:
	User: example@example.com
	Password copied to clipboard

	$ psst -r -a grandma_instagram
	Removed ~/.psst/.grandma_instagram
```