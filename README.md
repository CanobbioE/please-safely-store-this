# Please Safely Store This (PSST)
A secure command-line password manager written in Go.

(Yes, I did this only because I came up with the fun acronym)

## Key Features
- Secure local password storage
- AES-256 encryption
- Command-line interface for easy access
- Password generation
- Organization with tags
- Clipboard integration

## Installation

### Binaries
Binaries are provided [here](https://github.com/CanobbioE/please-safely-store-this/releases).

### From Source
1. Clone the repository
   ```
   git clone https://github.com/yourusername/passmanager.git
   cd passmanager
   ```

2. Build the project
   ```
   go build -o passmanager
   ```

   Or use the build script to build for multiple platforms:
   ```
   chmod +x build.sh
   ./build.sh
   ```

3. Install (optional)
   ```
   go install
   ```

## Usage

```
Usage:
  psst [command]

Available Commands:
  add         Add a new password entry
  completion  Generate the autocompletion script for the specified shell
  delete      Delete a password entry
  get         Retrieve a password
  help        Help about any command
  init        Initialize the password vault
  list        List all password entries
  update      Update an existing password

Flags:
      --config string   config file (default is $HOME/.psst/config.yaml)
  -h, --help            help for psst
  -v, --verbose         enable verbose output

```

### Common Use Cases

**Password Storage**
```
  psst add --service github --username dev@example.com --tags work,dev
  Enter password (or generate with -g): 
  Password stored in encrypted local vault.
```

**Password Retrieval**
```
  psst get gmail
```

**Update Existing Password**
```
  psst update gmail --password newSecurePass456
```

**List All Services**
```
  psst list
```

**Search for Specific Services**
```
  psst search google
```

**Copy Password to Clipboard**
```
  psst copy gmail
  Password copied to clipboard. Will clear in 30 seconds.
```

**Generate Strong Password**
```
  psst generate --length 16 --special-chars
```

**Export/Backup Passwords**
```
  psst export --file backup.enc
```

**Backup storage**
```
  psst backup --location ~/backups/
  Encrypted backup created at ~/backups/psst_backup_2025-04-14.enc
```

**Customize storage location**
```
  config set storage.path /path/to/custom/location/vault.db
  Storage location updated. Existing data will be migrated.
```

## Security Notes

- Passwords never leave your machine except for explicit exports
- The master password is never stored, only its hash
- Memory is securely wiped after use

## Roadmap
Many new features are yet to come, check our [roadmap](roadmap.md).

## License
[MIT](LICENSE)