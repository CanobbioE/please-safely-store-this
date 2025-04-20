package vault_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/CanobbioE/please-safely-store-this/internal/pkg/model"
	mockdb "github.com/CanobbioE/please-safely-store-this/internal/pkg/test/db"
	"github.com/CanobbioE/please-safely-store-this/internal/pkg/vault"
)

func TestNewManager(t *testing.T) {
	m := vault.NewManager(nil)
	if m == nil {
		t.Fatal("Manager is nil")
	}
}

func TestManager_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVault := mockdb.NewMockVault(ctrl)

	tests := []struct {
		name   string
		vault  vault.Vault
		preRun func(*testing.T)
	}{
		{
			name:  "successfully closes the vault",
			vault: mockVault,
			preRun: func(_ *testing.T) {
				mockVault.EXPECT().Close().Return(nil)
			},
		},
		{
			name:  "successfully closes the vault even if db throws an error",
			vault: mockVault,
			preRun: func(_ *testing.T) {
				mockVault.EXPECT().Close().Return(errors.New("db error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.preRun != nil {
				tt.preRun(t)
			}
			m := vault.NewManager(tt.vault)
			m.Close()
		})
	}
}

func TestManager_Unlock(t *testing.T) {
	var (
		mockVault *mockdb.MockVault
		manager   *vault.Manager
	)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	beforeEach := func(_ *testing.T) {
		mockVault = mockdb.NewMockVault(ctrl)
		manager = vault.NewManager(mockVault)
	}

	tests := []struct {
		name       string
		hashString string
	}{
		{
			name:       "too many fields",
			hashString: "$argon2id$v=19$m=65536,t=3,p=4$salt$hash$and$extra$fields",
		},
		{
			name:       "not argon2id",
			hashString: "$argon3.0d$v=19$m=65536,t=3,p=4$salt$hash$and$extra$fields",
		},
		{
			name:       "cannot read versions",
			hashString: "$argon2id$version19$m=65536,t=3,p=4$salt$hash$and$extra$fields",
		},
		{
			name:       "invalid salt",
			hashString: "$argon2id$v=19$m=65536,t=3,p=4$NOTASALT???!$hash$and$extra$fields",
		},
		{
			name:       "invalid hash",
			hashString: "$argon2id$v=19$m=65536,t=3,p=4$6c89d7fbb2e90bbe9e91509fc4d5b546$NOTANHASH??!!",
		},
	}

	for _, tt := range tests {
		t.Run(`returns error when it fails to split hash (`+tt.name+`)`, func(t *testing.T) {
			beforeEach(t)
			mockVault.EXPECT().
				GetVaultMetadata().
				Return(&model.VaultMetadata{
					MasterHash: tt.hashString,
				}, nil)
			unlocked, err := manager.Unlock("password")
			if unlocked {
				t.Fatal("Expected vault to be locked")
			}
			if err == nil {
				t.Fatal("Expected error, got nil")
			}
			if !strings.Contains(err.Error(), "invalid master hash format") {
				t.Fatalf("Expected error to contain 'invalid master hash format', got [%v]", err)
			}
		})
	}

	t.Run("returns error when it fails to get metadata", func(t *testing.T) {
		beforeEach(t)
		mockVault.EXPECT().GetVaultMetadata().Return(nil, errors.New("db error"))
		unlocked, err := manager.Unlock("password")
		if unlocked {
			t.Fatal("Expected vault to be locked")
		}

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if !strings.Contains(err.Error(), "failed to get vault metadata:") {
			t.Fatalf("Expected error to contain 'failed to get vault metadata:', got [%v]", err)
		}
	})

	t.Run("does not unlock with wrong password", func(t *testing.T) {
		beforeEach(t)
		mockVault.EXPECT().
			GetVaultMetadata().
			Return(&model.VaultMetadata{
				MasterHash: "$argon2id$v=19$m=65536,t=3,p=4$6c89d7fbb2e90bbe9e91509fc4d5b546$67b53292ba1f9c1c6c9193c48404d8c9fdfeb93041d5affcd08181241e284cdd", //nolint:lll
			}, nil)

		unlocked, err := manager.Unlock("incorrect password")
		if unlocked {
			t.Fatal("Expected vault to be locked")
		}

		if err != nil {
			t.Fatal("Unexpected error: ", err)
		}
	})

	t.Run("returns error when it fails to update metadata", func(t *testing.T) {
		beforeEach(t)
		meta := &model.VaultMetadata{
			MasterHash: "$argon2id$v=19$m=65536,t=3,p=4$6c89d7fbb2e90bbe9e91509fc4d5b546$67b53292ba1f9c1c6c9193c48404d8c9fdfeb93041d5affcd08181241e284cdd", //nolint:lll
			CreatedAt:  time.Now().Add(-24 * time.Hour).UTC(),
			LastAccess: time.Now().UTC(),
			Version:    "1.0.0",
		}
		mockVault.EXPECT().
			GetVaultMetadata().
			Return(meta, nil)

		mockVault.EXPECT().
			SaveVaultMetadata(gomock.Eq(meta)).
			Return(errors.New("db error"))

		unlocked, err := manager.Unlock("password123456")
		if !unlocked {
			t.Fatal("Expected vault to be unlocked")
		}

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if !strings.Contains(err.Error(), "failed to update last access time:") {
			t.Fatalf("Expected error to contain 'failed to update last access time:', got [%v]", err)
		}
	})

	t.Run("unlocks with correct password", func(t *testing.T) {
		beforeEach(t)
		meta := &model.VaultMetadata{
			MasterHash: "$argon2id$v=19$m=65536,t=3,p=4$6c89d7fbb2e90bbe9e91509fc4d5b546$67b53292ba1f9c1c6c9193c48404d8c9fdfeb93041d5affcd08181241e284cdd", //nolint:lll
			CreatedAt:  time.Now().Add(-24 * time.Hour).UTC(),
			LastAccess: time.Now().UTC(),
			Version:    "1.0.0",
		}
		mockVault.EXPECT().
			GetVaultMetadata().
			Return(meta, nil)

		mockVault.EXPECT().
			SaveVaultMetadata(gomock.Eq(meta)).
			Return(nil)

		unlocked, err := manager.Unlock("password123456")
		if !unlocked {
			t.Fatal("Expected vault to be unlocked")
		}

		if err != nil {
			t.Fatal("Unexpected error: ", err)
		}
	})

	t.Run("unlocks if already unlocked", func(t *testing.T) {
		beforeEach(t)
		meta := &model.VaultMetadata{
			MasterHash: "$argon2id$v=19$m=65536,t=3,p=4$6c89d7fbb2e90bbe9e91509fc4d5b546$67b53292ba1f9c1c6c9193c48404d8c9fdfeb93041d5affcd08181241e284cdd", //nolint:lll
			CreatedAt:  time.Now().Add(-24 * time.Hour).UTC(),
			LastAccess: time.Now().UTC(),
			Version:    "1.0.0",
		}
		mockVault.EXPECT().
			GetVaultMetadata().
			Return(meta, nil)

		mockVault.EXPECT().
			SaveVaultMetadata(gomock.Eq(meta)).
			Return(nil)

		unlocked, err := manager.Unlock("password123456")
		if !unlocked {
			t.Fatal("Expected vault to be unlocked")
		}

		if err != nil {
			t.Fatal("Unexpected error: ", err)
		}
		unlocked, err = manager.Unlock("password123456")
		if !unlocked {
			t.Fatal("Expected vault to be unlocked")
		}

		if err != nil {
			t.Fatal("Unexpected error: ", err)
		}
	})
}

// TODO: more tests
