package generator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/Andrei-hub11/archforge/internal/config"
	"github.com/Andrei-hub11/archforge/internal/templates"
)

// TestTemplateGeneration_CleanArchKeycloakPgEf tests the generation of the Clean Architecture with EF template
// and verifies that it can be built with `dotnet build`
func TestTemplateGeneration_CleanArchKeycloakPgEf(t *testing.T) {
	// Skip this test if we're not in a full test environment with dotnet installed
	if os.Getenv("INTEGRATION_TESTS") != "1" {
		t.Skip("Skipping integration test; set INTEGRATION_TESTS=1 to enable")
	}

	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "archforge-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set the templates directory for testing
	templates.SetTemplatesRootDir("../../templates") // Adjust this path as needed
	defer templates.ResetTemplatesRootDir()

	// Configure the project
	projectName := "TestEfProject"
	projectDir := filepath.Join(tempDir, projectName)
	cfg := config.ProjectConfig{
		Name:      projectName,
		OutputDir: tempDir,
		Template:  "clean-arch-keycloak-pg-ef",
	}

	// Generate the project
	err = Generate(cfg)
	if err != nil {
		t.Fatalf("Failed to generate project: %v", err)
	}

	// Verify the main project folders exist
	requiredDirs := []string{
		"src",
		"tests",
		filepath.Join("src", fmt.Sprintf("%s.Api", projectName)),
		filepath.Join("src", fmt.Sprintf("%s.Application", projectName)),
		filepath.Join("src", fmt.Sprintf("%s.Domain", projectName)),
		filepath.Join("src", fmt.Sprintf("%s.Infrastructure", projectName)),
	}

	for _, dir := range requiredDirs {
		fullPath := filepath.Join(projectDir, dir)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("Expected directory %s does not exist", fullPath)
		}
	}

	// Run dotnet build to verify the project can build
	cmd := exec.Command("dotnet", "build")
	cmd.Dir = projectDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build project: %v\nOutput: %s", err, output)
	}

	t.Logf("Successfully built project %s", projectName)
}

// TestTemplateGeneration_CleanArchKeycloakPgDapper tests the generation of the Clean Architecture with Dapper template
// and verifies that it can be built with `dotnet build`
func TestTemplateGeneration_CleanArchKeycloakPgDapper(t *testing.T) {
	// Skip this test if we're not in a full test environment with dotnet installed
	if os.Getenv("INTEGRATION_TESTS") != "1" {
		t.Skip("Skipping integration test; set INTEGRATION_TESTS=1 to enable")
	}

	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "archforge-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set the templates directory for testing
	templates.SetTemplatesRootDir("../../templates") // Adjust this path as needed
	defer templates.ResetTemplatesRootDir()

	// Configure the project
	projectName := "TestDapperProject"
	projectDir := filepath.Join(tempDir, projectName)
	cfg := config.ProjectConfig{
		Name:      projectName,
		OutputDir: tempDir,
		Template:  "clean-arch-keycloak-pg-dapper",
	}

	// Generate the project
	err = Generate(cfg)
	if err != nil {
		t.Fatalf("Failed to generate project: %v", err)
	}

	// Verify the main project folders exist
	requiredDirs := []string{
		"src",
		"tests",
		filepath.Join("src", fmt.Sprintf("%s.Api", projectName)),
		filepath.Join("src", fmt.Sprintf("%s.Application", projectName)),
		filepath.Join("src", fmt.Sprintf("%s.Domain", projectName)),
		filepath.Join("src", fmt.Sprintf("%s.Infrastructure", projectName)),
	}

	for _, dir := range requiredDirs {
		fullPath := filepath.Join(projectDir, dir)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("Expected directory %s does not exist", fullPath)
		}
	}

	// Run dotnet build to verify the project can build
	cmd := exec.Command("dotnet", "build")
	cmd.Dir = projectDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build project: %v\nOutput: %s", err, output)
	}

	t.Logf("Successfully built project %s", projectName)
}
