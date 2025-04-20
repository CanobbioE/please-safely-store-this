package vault

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/argon2"
)

// Hash parts.
type hashParts struct {
	algorithm string
	salt      []byte
	hash      []byte
	version   int
	memory    uint32
	time      uint32
	threads   uint8
}

// hashPassword hashes a password using Argon2id.
func hashPassword(password string, salt []byte) (hash string, key []byte) {
	// Argon2id parameters
	memory := uint32(64 * 1024) // 64MB
	time := uint32(3)           // 3 iterations
	threads := uint8(4)         // 4 threads
	keyLen := uint32(32)        // 32 bytes (256 bits)

	key = argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)
	hash = fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, memory, time, threads,
		hex.EncodeToString(salt),
		hex.EncodeToString(key))

	return hash, key
}

// verifyPassword verifies a password against a hash.
func verifyPassword(password, encodedHash string) bool {
	// Parse the hash
	parts := splitHash(encodedHash)
	if parts == nil {
		return false
	}

	// Hash the password with the same parameters
	_, key := hashPassword(password, parts.salt)

	// Compare the hashes
	return sha256.Sum256(key) == sha256.Sum256(parts.hash)
}

// splitHash splits an encoded hash into its parts.
func splitHash(encodedHash string) *hashParts {
	// Example: $argon2id$v=19$m=65536,t=3,p=4$salt$hash
	fields := splitString(encodedHash, '$')
	if len(fields) != 6 {
		return nil
	}

	if fields[1] != "argon2id" {
		return nil
	}

	var version int
	_, err := fmt.Sscanf(fields[2], "v=%d", &version)
	if err != nil {
		return nil
	}

	var memory, time uint32
	var threads uint8
	_, err = fmt.Sscanf(fields[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
	if err != nil {
		return nil
	}

	salt, err := hex.DecodeString(fields[4])
	if err != nil {
		return nil
	}

	hash, err := hex.DecodeString(fields[5])
	if err != nil {
		return nil
	}

	return &hashParts{
		algorithm: "argon2id",
		version:   version,
		memory:    memory,
		time:      time,
		threads:   threads,
		salt:      salt,
		hash:      hash,
	}
}

// splitString splits a string by a separator but without empty elements.
func splitString(s string, sep byte) []string {
	parts := make([]string, 0, 8)
	start := 0
	for i := range len(s) {
		if s[i] == sep {
			if i > start {
				parts = append(parts, s[start:i])
			}
			start = i + 1
		}
	}
	if start < len(s) {
		parts = append(parts, s[start:])
	}
	return parts
}
