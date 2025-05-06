package templates

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Andrei-hub11/archforge/internal/config"
)

// Para permitir configuração durante testes
var templatesRootDir string

type TemplateData struct {
	ProjectName string
}

// ProcessTemplateFile reads a template file, applies the data, and writes it to the destination
func ProcessTemplateFile(templatePath, destPath string, cfg config.ProjectConfig) error {
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo de template: %w", err)
	}

	if strings.TrimSuffix(destPath, ".tmpl") != destPath {
		destPath = strings.TrimSuffix(destPath, ".tmpl")
	}

	// if the content contains placeholder, replace it with the data
	if strings.Contains(string(content), "{{.ProjectName}}") {
		content = []byte(strings.ReplaceAll(string(content), "{{.ProjectName}}", cfg.Name))
	}

	if strings.Contains(string(content), "__ProjectName__") {
		content = []byte(strings.ReplaceAll(string(content), "__ProjectName__", cfg.Name))
	}

	tmpl, err := template.New(filepath.Base(templatePath)).Parse(string(content))
	if err != nil {
		return fmt.Errorf("erro ao analisar template: %w", err)
	}

	// Ensure the destination directory exists
	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório de destino: %w", err)
	}

	file, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo de destino: %w", err)
	}

	defer file.Close()

	// setting the data for the template
	data := TemplateData{
		ProjectName: cfg.Name,
	}

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("erro ao aplicar template: %w", err)
	}

	return nil
}

// CopyTemplateDir recursively copies templates from a source directory to a destination
// and processes template files with the provided data
func CopyTemplateDir(sourceDir, destDir string, cfg config.ProjectConfig) error {
	// Get list of template files
	entries, err := os.ReadDir(sourceDir)
	if err != nil {
		return fmt.Errorf("erro ao ler diretório de templates: %w", err)
	}

	for _, entry := range entries {

		// Ignore obj and bin directories
		if entry.IsDir() && (entry.Name() == "obj" || entry.Name() == "bin") {
			continue
		}

		sourcePath := filepath.Join(sourceDir, entry.Name())
		destPath := filepath.Join(destDir, entry.Name())

		// If the filename contains a template placeholder, replace it
		if strings.Contains(entry.Name(), "{{.ProjectName}}") {
			destPath = filepath.Join(destDir,
				strings.ReplaceAll(entry.Name(), "{{.ProjectName}}", cfg.Name))
		}

		if entry.IsDir() {
			// Create destination directory
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return fmt.Errorf("erro ao criar diretório: %w", err)
			}
			// Recursive call for subdirectory
			if err := CopyTemplateDir(sourcePath, destPath, cfg); err != nil {
				return err
			}
		} else {
			// Process file
			if err := ProcessTemplateFile(sourcePath, destPath, cfg); err != nil {
				return err
			}
		}
	}

	return nil
}

// GetTemplatesRootDir returns the path to the templates directory
func GetTemplatesRootDir() string {
	// If the templates directory has been explicitly configured, use it
	if templatesRootDir != "" {
		return templatesRootDir
	}

	// Get the executable path
	exePath, err := os.Executable()
	if err != nil {
		// Fallback to current directory if executable path cannot be determined
		return "templates"
	}

	// The templates directory is expected to be in the same directory as the executable
	return filepath.Join(filepath.Dir(exePath), "templates")
}

// SetTemplatesRootDir permite definir o diretório raiz de templates
// Útil para testes ou personalização
func SetTemplatesRootDir(dir string) {
	templatesRootDir = dir
}

// ResetTemplatesRootDir restaura o diretório raiz de templates para o padrão
func ResetTemplatesRootDir() {
	templatesRootDir = ""
}

// GenerateWebApi creates a Web API application from templates
func GenerateWebApi(projectDir string, cfg config.ProjectConfig) error {
	templateDir := filepath.Join(GetTemplatesRootDir(), "webapi")
	return CopyTemplateDir(templateDir, projectDir, cfg)
}

func GenerateCleanArchKeycloakPgDapper(projectDir string, cfg config.ProjectConfig) error {
	templateDir := filepath.Join(GetTemplatesRootDir(), "clean-arch-keycloak-pg-dapper")
	return CopyTemplateDir(templateDir, projectDir, cfg)
}

func GenerateCleanArchKeycloakPgEf(projectDir string, cfg config.ProjectConfig) error {
	templateDir := filepath.Join(GetTemplatesRootDir(), "clean-arch-keycloak-pg-ef")
	return CopyTemplateDir(templateDir, projectDir, cfg)
}

func PrintBuildTree(directoryPath string, prefix string) {
	entries, err := os.ReadDir(directoryPath)
	if err != nil {
		fmt.Println("erro ao ler diretório de templates: %w", err)
		return
	}

	for i, entry := range entries {
		// Ignore obj and bin directories
		if entry.IsDir() && (entry.Name() == "obj" || entry.Name() == "bin") {
			continue
		}

		connector := "├──"
		if i == len(entries)-1 {
			connector = "└──"
		}

		fmt.Printf("%s%s %s\n", prefix, connector, entry.Name())

		if entry.IsDir() {
			newPrefix := prefix

			if i == len(entries)-1 {
				newPrefix += "    "
			} else {
				newPrefix += "│   "
			}

			PrintBuildTree(filepath.Join(directoryPath, entry.Name()), newPrefix)
		}
	}
}
