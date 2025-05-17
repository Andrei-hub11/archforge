# ArchForge

ArchForge is a command-line tool for generating the basic structure of a project using templates.  
You can quickly create a new project by specifying the name, path, and desired template.

## Usage

ArchForge offers two ways to create projects:

### Interactive Mode

Run the interactive mode which will guide you through the project creation:

```bash
./archforge interactive
```

The interactive mode will ask for:

1. Project name (default: MyProject)
2. Output directory (default: current directory)
3. Template to use (select from available options)
4. Preview option (Yes/No)

### Command Line Mode

Create a project directly using command line arguments:

```bash
./archforge create --name MyProject --template webapi [--output ./path/to/output] [--preview]
```

Required flags:

- `--name` or `-n`: The name of your project (must not contain spaces or hyphens)
- `--template` or `-t`: The template to use

Optional flags:

- `--output` or `-o`: Output directory (defaults to current directory)
- `--preview` or `-p`: Preview the project structure before generating

Available templates:

- `webapi`: Basic ASP.NET Core Web API
- `clean-arch-keycloak-pg-dapper`: Clean Architecture with Keycloak, PostgreSQL, and Dapper
- `clean-arch-keycloak-pg-ef`: Clean Architecture with Keycloak, PostgreSQL, and Entity Framework

## Testing

To run the tests, make sure you have Go installed and follow these steps:

1. Clone the repository:

```bash
git clone https://github.com/Andrei-hub11/archforge.git
cd archforge
```

2. Run all tests:

```bash
go test ./...
```

3. Run tests for a specific package:

```bash
go test ./cmd
```

4. Run tests with verbose output:

```bash
go test -v ./...
```

Note: Some interactive tests may be skipped in CI environments as they require user input. These tests are automatically skipped when the `CI` environment variable is set.

## Integration Tests

ArchForge includes integration tests that verify the template generation and compilation. These tests:

1. Generate projects from each template
2. Verify the correct folder structure and file creation
3. Test that the generated projects can be built with `dotnet build`
4. Verify that the basic WebAPI template can be run with `dotnet run`

### Running Integration Tests

Integration tests require that you have the .NET SDK installed on your system.

#### Windows

Run the PowerShell script:

```powershell
./run_integration_tests.ps1
```

#### Linux/macOS

Run the bash script:

```bash
chmod +x ./run_integration_tests.sh
./run_integration_tests.sh
```

Alternatively, you can run specific integration tests:

```bash
# Set environment variable
export INTEGRATION_TESTS=1  # or $env:INTEGRATION_TESTS=1 in PowerShell

# Template generation tests
go test -v ./internal/generator -run TestTemplateGeneration_

# Folder structure tests
go test -v ./internal/generator -run TestFolderStructure_

# Project execution tests
go test -v ./internal/generator -run TestProjectExecution_
```

These tests verify that templates can be properly generated, that they have the correct structure, and that the generated code can be built and run successfully.
