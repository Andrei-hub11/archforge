package cmd

import (
	"strings"
	"testing"
)

// TestCreateCommand_RequiredFlags tests that required flags are properly enforced
func TestCreateCommand_RequiredFlags(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantErr     bool
		errorString string
	}{
		{
			name:        "missing name flag",
			args:        []string{"create", "--template", "webapi"},
			wantErr:     true,
			errorString: "required flag(s) \"name\" not set",
		},
		{
			name:        "missing template flag",
			args:        []string{"create", "--name", "myproject"},
			wantErr:     true,
			errorString: "required flag(s) \"template\" not set",
		},
		{
			name:        "missing all required flags",
			args:        []string{"create"},
			wantErr:     true,
			errorString: "required flag(s)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootCmd := CreateRootCommand()

			output, err := ExecuteCommand(rootCmd, tt.args...)

			if (err != nil) != tt.wantErr {
				t.Errorf("Expected error: %v, got: %v", tt.wantErr, err != nil)
				return
			}

			if tt.wantErr && !strings.Contains(output, tt.errorString) {
				t.Errorf("Expected error string containing '%s', got: '%s'", tt.errorString, output)
			}
		})
	}
}

// TestValidateProjectName tests the project name validation function
func TestValidateProjectName(t *testing.T) {
	tests := []struct {
		name        string
		projectName string
		wantErr     bool
		errorMsg    string
	}{
		{
			name:        "valid name",
			projectName: "ValidProject",
			wantErr:     false,
		},
		{
			name:        "empty name",
			projectName: "",
			wantErr:     true,
			errorMsg:    "project name is required",
		},
		{
			name:        "name with spaces",
			projectName: "Project With Spaces",
			wantErr:     true,
			errorMsg:    "project name must not contain spaces",
		},
		{
			name:        "name with hyphens",
			projectName: "project-with-hyphens",
			wantErr:     true,
			errorMsg:    "project name must not contain hyphens",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateProjectName(tt.projectName)

			if (err != nil) != tt.wantErr {
				t.Errorf("validateProjectName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && !strings.Contains(err.Error(), tt.errorMsg) {
				t.Errorf("validateProjectName() error message = %v, want to contain %v", err.Error(), tt.errorMsg)
			}
		})
	}
}

// CreateRootCommand creates a root command with the create command as a subcommand
func CreateRootCommand() *Command {
	rootCmd := &Command{Use: "root"}
	rootCmd.AddCommand(CreateCmd())
	return rootCmd
}
