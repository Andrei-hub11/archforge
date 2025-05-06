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
	// Verifica se o diretório de saída existe
	projectDir := filepath.Join(cfg.OutputDir, cfg.Name)
	if _, err := os.Stat(projectDir); !os.IsNotExist(err) {
		return fmt.Errorf("o diretório '%s' já existe", projectDir)
	}

	// Cria o diretório do projeto
	err := os.MkdirAll(projectDir, 0755)
	if err != nil {
		return fmt.Errorf("erro ao criar diretório do projeto: %w", err)
	}

	// Gera os arquivos apropriados com base no template selecionado
	switch cfg.Template {
	case "webapi":
		err = templates.GenerateWebApi(projectDir, cfg)
	case "clean-arch-keycloak-pg-dapper":
		err = templates.GenerateCleanArchKeycloakPgDapper(projectDir, cfg)
	case "clean-arch-keycloak-pg-ef":
		err = templates.GenerateCleanArchKeycloakPgEf(projectDir, cfg)
	default:
		return fmt.Errorf("template desconhecido: %s", cfg.Template)
	}

	if err != nil {
		return fmt.Errorf("erro ao gerar template: %w", err)
	}

	return nil
}

func GenerateBuildTree(cfg config.ProjectConfig) error {
	templateDir := filepath.Join(templates.GetTemplatesRootDir(), cfg.Template)
	templates.PrintBuildTree(templateDir, "")
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
