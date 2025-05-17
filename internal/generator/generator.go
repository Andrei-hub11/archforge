// projeto/internal/generator/generator.go
package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Andrei-hub11/archforge/internal/config"
	"github.com/Andrei-hub11/archforge/internal/templates"
)

// Generator interface allows mocking for testing
type Generator interface {
	Generate(cfg config.ProjectConfig) error
}

// DefaultGenerator is the real implementation
type DefaultGenerator struct{}

// Generate creates a new project based on the provided configuration
func (g *DefaultGenerator) Generate(cfg config.ProjectConfig) error {
	// Check if the output directory exists
	projectDir := filepath.Join(cfg.OutputDir, cfg.Name)
	if _, err := os.Stat(projectDir); !os.IsNotExist(err) {
		return fmt.Errorf("directory '%s' already exists", projectDir)
	}

	// Create the project directory
	err := os.MkdirAll(projectDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Generate appropriate files based on the selected template
	switch cfg.Template {
	case "webapi":
		err = templates.GenerateWebApi(projectDir, cfg)
	case "clean-arch-keycloak-pg-dapper":
		err = templates.GenerateCleanArchKeycloakPgDapper(projectDir, cfg)
	case "clean-arch-keycloak-pg-ef":
		err = templates.GenerateCleanArchKeycloakPgEf(projectDir, cfg)
	default:
		return fmt.Errorf("unknown template: %s", cfg.Template)
	}

	if err != nil {
		return fmt.Errorf("failed to generate template: %w", err)
	}

	return nil
}

func GenerateBuildTree(cfg config.ProjectConfig) error {
	templateDir := filepath.Join(templates.GetTemplatesRootDir(), cfg.Template)
	templates.PrintBuildTree(templateDir, "", cfg.Name)
	return nil
}

// Default generator instance
var defaultGenerator Generator = &DefaultGenerator{}

// Generate is a package-level function that uses the default generator
func Generate(cfg config.ProjectConfig) error {
	return defaultGenerator.Generate(cfg)
}

// For testing purposes
func SetGenerator(g Generator) {
	defaultGenerator = g
}

// Reset to default generator
func ResetGenerator() {
	defaultGenerator = &DefaultGenerator{}
}
