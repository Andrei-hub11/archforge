package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/Andrei-hub11/archforge/internal/config"
	"github.com/Andrei-hub11/archforge/internal/templates"
)

// Test folder structure for the clean architecture templates
func TestFolderStructure_CleanArchitectureTemplates(t *testing.T) {
	// Skip this test if we're not in a full test environment
	if os.Getenv("INTEGRATION_TESTS") != "1" {
		t.Skip("Skipping integration test; set INTEGRATION_TESTS=1 to enable")
	}

	templatesList := []string{
		"clean-arch-keycloak-pg-ef",
		"clean-arch-keycloak-pg-dapper",
	}

	for _, templateName := range templatesList {
		t.Run(templateName, func(t *testing.T) {
			// Create a temporary directory for the test
			tempDir, err := os.MkdirTemp("", fmt.Sprintf("archforge-folders-%s-*", templateName))
			if err != nil {
				t.Fatalf("Failed to create temp directory: %v", err)
			}
			defer os.RemoveAll(tempDir)

			// Set the templates directory for testing
			templates.SetTemplatesRootDir("../../templates") // Adjust this path as needed
			defer templates.ResetTemplatesRootDir()

			// Configure the project
			projectName := "TestFolders"
			projectDir := filepath.Join(tempDir, projectName)
			cfg := config.ProjectConfig{
				Name:      projectName,
				OutputDir: tempDir,
				Template:  templateName,
			}

			// Generate the project
			err = Generate(cfg)
			if err != nil {
				t.Fatalf("Failed to generate project: %v", err)
			}

			// Expected folder structure for clean architecture templates
			expectedFolders := []string{
				"src",
				"tests",
				filepath.Join("src", fmt.Sprintf("%s.Api", projectName)),
				filepath.Join("src", fmt.Sprintf("%s.Application", projectName)),
				filepath.Join("src", fmt.Sprintf("%s.Domain", projectName)),
				filepath.Join("src", fmt.Sprintf("%s.Infrastructure", projectName)),
			}

			// Expected files that should exist
			expectedFiles := []string{
				fmt.Sprintf("%s.sln", projectName),
				filepath.Join("src", fmt.Sprintf("%s.Api", projectName), "Program.cs"),
				filepath.Join("src", fmt.Sprintf("%s.Api", projectName), "appsettings.json"),
				filepath.Join("src", fmt.Sprintf("%s.Domain", projectName), "Entities", "Product.cs"),
				filepath.Join("src", fmt.Sprintf("%s.Domain", projectName), "Repositories", "IProductRepository.cs"),
				"realm-export.json",
				"README.md",
			}

			// Check if the expected folders exist
			for _, folder := range expectedFolders {
				folderPath := filepath.Join(projectDir, folder)
				if _, err := os.Stat(folderPath); os.IsNotExist(err) {
					t.Errorf("Expected folder %s does not exist", folderPath)
				}
			}

			// Check if the expected files exist
			for _, file := range expectedFiles {
				filePath := filepath.Join(projectDir, file)
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					t.Errorf("Expected file %s does not exist", filePath)
				}
			}

			t.Logf("Successfully verified folder structure for %s template", templateName)
		})
	}
}
