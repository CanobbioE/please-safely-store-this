package config_test

import (
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/CanobbioE/please-safely-store-this/internal/pkg/config"
)

func TestDefaultConfig(t *testing.T) {
	got := config.DefaultConfig()
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}

	matchConfig(t, got, &config.Config{
		AutoLockTimeout:     15 * time.Minute,
		ClipboardTimeout:    30 * time.Second,
		ShowPasswords:       false,
		BackupCount:         5,
		PasswordLength:      16,
		UseSpecialChars:     true,
		UseNumbers:          true,
		UseUppercase:        true,
		MinPasswordStrength: 2,
		DBPath:              filepath.Join(home, ".psst", "vault.db"),
		BackupDir:           filepath.Join(home, ".psst", "backups"),
	})
}

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    *config.Config
		postRun func(*testing.T)
	}{
		{
			name: "successfully loads config from file",
			path: "./testdata/config_test.yml",
			want: &config.Config{
				DBPath:              "mock_db_path",
				BackupDir:           "mock_backup_dir",
				AutoLockTimeout:     1 * time.Minute,
				ClipboardTimeout:    1 * time.Second,
				BackupCount:         42,
				PasswordLength:      42,
				MinPasswordStrength: 42,
				ShowPasswords:       true,
				UseSpecialChars:     true,
				UseNumbers:          true,
				UseUppercase:        true,
			},
		},
		{
			name: "successfully loads default config when config file does not exist",
			path: "non-existing-config.yml",
			want: config.DefaultConfig(),
			postRun: func(t *testing.T) {
				if err := os.Remove("./non-existing-config.yml"); err != nil {
					t.Errorf("Failed to remove non-existing config file")
				}
			},
		},
		{
			name: "successfully loads default config when config file cannot be parsed",
			path: "./testdata/unparsable.yml",
			want: config.DefaultConfig(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := config.LoadConfig(tt.path)
			matchConfig(t, got, tt.want)
			if tt.postRun != nil {
				tt.postRun(t)
			}
		})
	}
}

func Test_SaveConfig(t *testing.T) {
	tmpDir := t.TempDir()
	tests := []struct {
		name        string
		config      *config.Config
		path        string
		expectedErr string
	}{
		{
			name:   "successfully save config",
			config: config.DefaultConfig(),
			path:   path.Join(tmpDir, "tmp-psst-config.yml"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := config.SaveConfig(tt.config, tt.path)
			switch {
			case err == nil && tt.expectedErr != "":
				t.Fatalf("Expected error containing [%s], but got no error", tt.expectedErr)
			case err != nil && tt.expectedErr == "":
				t.Fatalf("Expected no error, got [%v]", err)
			case err != nil && !strings.Contains(err.Error(), tt.expectedErr):
				t.Fatalf("Expected error containing [%s], but got [%v]", tt.expectedErr, err)
			case tt.expectedErr == "":
				data, err := os.ReadFile(tt.path)
				if err != nil {
					t.Fatalf("Failed to open stored config: %v", err)
				}
				var got config.Config
				if err = yaml.Unmarshal(data, &got); err != nil {
					t.Fatalf("Failed to parse config file: %s", err)
				}
				matchConfig(t, tt.config, &got)
			}
		})
	}
}

func matchConfig(t *testing.T, got, want *config.Config, skipFields ...string) {
	t.Helper()

	skip := make(map[string]struct{}, len(skipFields))
	for _, skipField := range skipFields {
		skip[skipField] = struct{}{}
	}

	gotVal := reflect.Indirect(reflect.ValueOf(got))
	wantVal := reflect.Indirect(reflect.ValueOf(want))
	typ := gotVal.Type()

	for i := range gotVal.NumField() {
		fieldName := typ.Field(i).Name
		gotField := gotVal.Field(i).Interface()
		wantField := wantVal.Field(i).Interface()
		if _, shouldSkip := skip[fieldName]; shouldSkip {
			continue
		}

		if !reflect.DeepEqual(gotField, wantField) {
			t.Fatalf("expected %s to equal [%v], instead got [%v]", fieldName, wantField, gotField)
		}
	}
}
