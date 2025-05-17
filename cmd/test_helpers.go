package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

// Command is a shortcut for cobra.Command
type Command = cobra.Command

// ExecuteCommand executes a command with arguments and returns output and error
func ExecuteCommand(cmd *cobra.Command, args ...string) (string, error) {
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)

	err := cmd.Execute()
	return buf.String(), err
}

// TestingHelper provides utilities for testing
type TestingHelper struct {
	T            *testing.T
	TempDir      string
	OutputDir    string
	TemplatesDir string
}

// NewTestingHelper creates a new testing helper
func NewTestingHelper(t *testing.T) *TestingHelper {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "archforge-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Create output and templates subdirectories
	outputDir := filepath.Join(tempDir, "output")
	templatesDir := filepath.Join(tempDir, "templates")

	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	err = os.MkdirAll(templatesDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create templates directory: %v", err)
	}

	return &TestingHelper{
		T:            t,
		TempDir:      tempDir,
		OutputDir:    outputDir,
		TemplatesDir: templatesDir,
	}
}

// Cleanup removes temporary test directories
func (h *TestingHelper) Cleanup() {
	if h.TempDir != "" {
		if err := os.RemoveAll(h.TempDir); err != nil {
			h.T.Logf("Warning: Failed to clean up temporary directory %s: %v", h.TempDir, err)
		}
	}
}

// SetupTemplates creates sample templates for testing
func (h *TestingHelper) SetupTemplates() {
	// Template directories to create
	templates := []string{
		"webapi",
		"clean-arch-keycloak-pg-dapper",
		"clean-arch-keycloak-pg-ef",
	}

	for _, tmpl := range templates {
		// Create template directory
		templateDir := filepath.Join(h.TemplatesDir, tmpl)
		err := os.MkdirAll(templateDir, 0755)
		if err != nil {
			h.T.Fatalf("Failed to create template directory %s: %v", templateDir, err)
		}

		// Create a sample file in each template directory
		programFile := filepath.Join(templateDir, "Program.cs")
		err = os.WriteFile(programFile, []byte("// Sample Program.cs for template "+tmpl), 0644)
		if err != nil {
			h.T.Fatalf("Failed to create sample file in template %s: %v", tmpl, err)
		}
	}

	// Make sure the templates directory exists
	if _, err := os.Stat(h.TemplatesDir); os.IsNotExist(err) {
		h.T.Fatalf("Templates directory doesn't exist after creation: %s", h.TemplatesDir)
	}

	// Print templates directory for debugging
	h.T.Logf("Templates directory set to: %s", h.TemplatesDir)

	// Set the templates directory for the template package
	if err := os.Setenv("ARCHFORGE_TEMPLATES_DIR", h.TemplatesDir); err != nil {
		h.T.Fatalf("Failed to set environment variable ARCHFORGE_TEMPLATES_DIR: %v", err)
	}

	// Verify the environment variable was set
	if dir := os.Getenv("ARCHFORGE_TEMPLATES_DIR"); dir != h.TemplatesDir {
		h.T.Fatalf("Failed to set ARCHFORGE_TEMPLATES_DIR. Expected: %s, Got: %s", h.TemplatesDir, dir)
	}
}

// VerifyDirExists checks if a directory exists
func VerifyDirExists(t *testing.T, path string) {
	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		t.Errorf("Directory does not exist: %s", path)
		return
	}
	if err != nil {
		t.Errorf("Error checking directory: %s, %v", path, err)
		return
	}
	if !stat.IsDir() {
		t.Errorf("Path is not a directory: %s", path)
	}
}

// VerifyFileExists checks if a file exists
func VerifyFileExists(t *testing.T, path string) {
	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		t.Errorf("File does not exist: %s", path)
		return
	}
	if err != nil {
		t.Errorf("Error checking file: %s, %v", path, err)
		return
	}
	if stat.IsDir() {
		t.Errorf("Path is a directory, not a file: %s", path)
	}
}
