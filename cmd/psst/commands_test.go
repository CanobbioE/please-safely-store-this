package psst_test

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/CanobbioE/please-safely-store-this/cmd/psst"
)

// TODO: expand test cases.

func TestCommands(t *testing.T) {
	tests := []struct {
		name        string
		cmd         *cobra.Command
		args        []string
		expectedErr string
	}{
		{
			name: "AddCmd successfully adds a password",
			cmd:  psst.AddCmd(),
			args: []string{
				"--service", "gmail",
				"--username", "user@example.com",
				"--password", "secret123",
				"--tags", "email,important",
			},
		},
		{
			name: "AddCmd does nothing if service is empty",
			cmd:  psst.AddCmd(),
			args: []string{
				"--service", "",
				"--username", "user@example.com",
				"--password", "secret123",
				"--tags", "email,important",
			},
		},
		{
			name: "AddCmd fails if service is missing",
			cmd:  psst.AddCmd(),
			args: []string{
				"--username", "user@example.com",
				"--password", "secret123",
				"--tags", "email,important",
			},
			expectedErr: `required flag(s) "service" not set`,
		},
		{
			name: "AddCmd does nothing if password is missing",
			cmd:  psst.AddCmd(),
			args: []string{
				"--service", "gmail",
				"--username", "user@example.com",
				"--tags", "email,important",
			},
		},
		// get
		{
			name: "GetCmd successfully gets a password",
			cmd:  psst.GetCmd(),
			args: []string{"--service", "gmail"},
		},
		{
			name: "GetCmd does nothing if service is empty",
			cmd:  psst.GetCmd(),
			args: []string{"--service", ""},
		},
		{
			name:        "GetCmd fails if service is missing",
			cmd:         psst.GetCmd(),
			expectedErr: `required flag(s) "service" not set`,
		},
		// list
		{
			name: "ListCmd successfully list all passwords",
			cmd:  psst.ListCmd(),
		},
		// update
		{
			name: "UpdateCmd successfully updates a password",
			cmd:  psst.UpdateCmd(),
			args: []string{"--service", "gmail", "--password", "secret123"},
		},
		{
			name: "UpdateCmd does nothing if service is empty",
			cmd:  psst.UpdateCmd(),
			args: []string{"--service", ""},
		},
		{
			name:        "UpdateCmd fails if service is missing",
			cmd:         psst.UpdateCmd(),
			expectedErr: `required flag(s) "service" not set`,
		},
		{
			name: "UpdateCmd does nothing if password is missing",
			cmd:  psst.UpdateCmd(),
			args: []string{"--service", "gmail"},
		},
		// delete
		{
			name: "DeleteCmd successfully deletes a password",
			cmd:  psst.DeleteCmd(),
			args: []string{"--service", "gmail"},
		},
		{
			name: "DeleteCmd does nothing if service is empty",
			cmd:  psst.DeleteCmd(),
			args: []string{"--service", ""},
		},
		{
			name:        "DeleteCmd fails if service is missing",
			cmd:         psst.DeleteCmd(),
			expectedErr: `required flag(s) "service" not set`,
		},
		// init
		{
			name: "InitCmd successfully initialize the password vault",
			cmd:  psst.InitCmd(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cmd.SetArgs(tt.args)
			err := tt.cmd.Execute()
			switch {
			case err == nil && tt.expectedErr != "":
				t.Fatalf("Expected error containing [%s], but got no error", tt.expectedErr)
			case err != nil && tt.expectedErr == "":
				t.Fatalf("Expected no error, got [%v]", err)
			case err != nil && !strings.Contains(err.Error(), tt.expectedErr):
				t.Fatalf("Expected error containing [%s], but got [%v]", tt.expectedErr, err)
			}
		})
	}
}
