package vault_test

import (
	"testing"

	"github.com/CanobbioE/please-safely-store-this/internal/pkg/vault"
)

func TestNewManager(t *testing.T) {
	m, err := vault.NewManager(t.TempDir() + "/vault.db")
	if err != nil {
		t.Fatal(err)
	}
	if m == nil {
		t.Fatal("Manager is nil")
	}
}

// TODO: more tests
