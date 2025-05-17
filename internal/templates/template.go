package templates

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Andrei-hub11/archforge/internal/config"
)

// Templates root directory for testing configuration
var templatesRootDir string

type TemplateData struct {
	ProjectName string
}

// ProcessTemplateFile reads a template file, applies the data, and writes it to the destination
func ProcessTemplateFile(templatePath, destPath string, cfg config.ProjectConfig) error {
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template file: %w", err)
	}

	destPath = ReplaceTemplatePlaceholderFromFileName(destPath, cfg.Name)

	// if the content contains placeholder, replace it with the data
	if strings.Contains(string(content), "{{.ProjectName}}") {
		content = []byte(strings.ReplaceAll(string(content), "{{.ProjectName}}", cfg.Name))
	}

	if strings.Contains(string(content), "__ProjectName__") {
		content = []byte(strings.ReplaceAll(string(content), "__ProjectName__", cfg.Name))
	}

	tmpl, err := template.New(filepath.Base(templatePath)).Parse(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Ensure the destination directory exists
	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	file, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("failed to close file: %w", closeErr)
		}
	}()

	// setting the data for the template
	data := TemplateData{
		ProjectName: cfg.Name,
	}

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to apply template: %w", err)
	}

	return nil
}

func ReplaceTemplatePlaceholderFromFolderName(folderName string, projectName string) string {
	if strings.Contains(folderName, "{{.ProjectName}}") {
		folderName = strings.ReplaceAll(folderName, "{{.ProjectName}}", projectName)
	}

	return folderName
}

// ReplaceTemplatePlaceholderFromFileName replaces the .tmpl or {{.ProjectName}} suffix from the file name
func ReplaceTemplatePlaceholderFromFileName(fileName string, projectName string) string {
	if strings.TrimSuffix(fileName, ".tmpl") != fileName {
		fileName = strings.TrimSuffix(fileName, ".tmpl")
	}

	if strings.Contains(fileName, "{{.ProjectName}}") {
		fileName = strings.ReplaceAll(fileName, "{{.ProjectName}}", projectName)
	}

	return fileName
}

// CopyTemplateDir recursively copies templates from a source directory to a destination
// and processes template files with the provided data
func CopyTemplateDir(sourceDir, destDir string, cfg config.ProjectConfig) error {
	// Get list of template files
	entries, err := os.ReadDir(sourceDir)
	if err != nil {
		return fmt.Errorf("failed to read templates directory: %w", err)
	}

	for _, entry := range entries {
		// Skip .vs directory and obj/bin directories
		if entry.IsDir() && (entry.Name() == ".vs" || entry.Name() == "obj" || entry.Name() == "bin") {
			continue
		}

		sourcePath := filepath.Join(sourceDir, entry.Name())
		destPath := filepath.Join(destDir, entry.Name())

		if entry.IsDir() {
			// Create destination directory
			destDirName := entry.Name()
			if strings.Contains(destDirName, "{{.ProjectName}}") {
				destDirName = strings.ReplaceAll(destDirName, "{{.ProjectName}}", cfg.Name)
			}
			destPath = filepath.Join(destDir, destDirName)

			if err := os.MkdirAll(destPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}

			// Recursive call for subdirectory
			if err := CopyTemplateDir(sourcePath, destPath, cfg); err != nil {
				return err
			}
		} else {
			// Skip binary files and VS-specific files
			if isBinaryOrVSFile(entry.Name()) {
				// Just copy the file without processing
				if err := copyFile(sourcePath, destPath); err != nil {
					return fmt.Errorf("failed to copy binary file: %w", err)
				}
				continue
			}

			// Process template file
			if err := ProcessTemplateFile(sourcePath, destPath, cfg); err != nil {
				return err
			}
		}
	}

	return nil
}

// isBinaryOrVSFile checks if a file is likely to be binary or VS-specific
func isBinaryOrVSFile(fileName string) bool {
	// List of extensions and patterns that indicate binary or VS-specific files
	binaryExtensions := []string{
		".bin", ".exe", ".dll", ".pdb", ".cache",
		".suo", ".user", ".vsidx", ".testlog",
		".v2", ".dtbcache", ".manifest"}

	ext := strings.ToLower(filepath.Ext(fileName))
	for _, binaryExt := range binaryExtensions {
		if ext == binaryExt {
			return true
		}
	}
	return false
}

// copyFile copies a file from src to dst without processing it as a template
func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(dst), 0755)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, input, 0644)
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

// SetTemplatesRootDir sets the templates root directory
// Useful for testing or customization
func SetTemplatesRootDir(dir string) {
	templatesRootDir = dir
}

// ResetTemplatesRootDir resets the templates root directory to the default
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

func PrintBuildTree(directoryPath string, prefix string, params ...string) {
	entries, err := os.ReadDir(directoryPath)
	if err != nil {
		fmt.Println("failed to read templates directory:", err)
		return
	}

	for i, entry := range entries {
		// Ignore obj and bin directories
		if entry.IsDir() && (entry.Name() == "obj" || entry.Name() == "bin" || entry.Name() == ".vs") {
			continue
		}

		connector := "├──"
		if i == len(entries)-1 {
			connector = "└──"
		}

		var entryName string

		if entry.IsDir() {
			entryName = ReplaceTemplatePlaceholderFromFolderName(entry.Name(), params[0])
		} else {
			entryName = ReplaceTemplatePlaceholderFromFileName(entry.Name(), params[0])
		}

		fmt.Printf("%s%s %s\n", prefix, connector, entryName)

		if entry.IsDir() {
			newPrefix := prefix

			if i == len(entries)-1 {
				newPrefix += "    "
			} else {
				newPrefix += "│   "
			}

			PrintBuildTree(filepath.Join(directoryPath, entry.Name()), newPrefix, params[0])
		}
	}
}
