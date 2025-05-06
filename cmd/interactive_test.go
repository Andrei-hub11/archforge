package cmd

import (
	"os"
	"testing"
)

// TestInteractiveCommand_Basic tests the basic functionality of the interactive command
func TestInteractiveCommand_Basic(t *testing.T) {
	// Skip this test in automated environments since it requires user input
	if os.Getenv("CI") != "" {
		t.Skip("Skipping interactive tests in CI environment")
	}

	// This test requires manual interaction, so we'll just check that the command
	// can be created and registered without errors
	rootCmd := CreateRootCommand()

	// Add the interactive command to the root command
	rootCmd.AddCommand(InteractiveCmd())

	// Verify the command is registered
	cmd, _, err := rootCmd.Find([]string{"interactive"})
	if err != nil {
		t.Errorf("Expected interactive command to be registered, got error: %v", err)
	}

	if cmd.Name() != "interactive" {
		t.Errorf("Expected command name to be 'interactive', got: %s", cmd.Name())
	}
}

// TestInteractiveCommand_ProjectNameValidation tests the project name validation
// in the interactive command
func TestInteractiveCommand_ProjectNameValidation(t *testing.T) {
	tests := []struct {
		name        string
		projectName string
		wantErr     bool
	}{
		{
			name:        "valid name",
			projectName: "ValidProject",
			wantErr:     false,
		},
		{
			name:        "name with spaces",
			projectName: "Project With Spaces",
			wantErr:     true,
		},
		{
			name:        "name with hyphens",
			projectName: "project-with-hyphens",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateProjectName(tt.projectName)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateProjectName(%s) error = %v, wantErr %v",
					tt.projectName, err, tt.wantErr)
			}
		})
	}
}
