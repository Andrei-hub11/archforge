# PowerShell script to run integration tests for Archforge
# This script will set the INTEGRATION_TESTS environment variable and run the Go tests

Write-Host "Running Archforge integration tests..." -ForegroundColor Green

try {
    $dotnetVersion = (dotnet --version)
    Write-Host "Found .NET SDK version: $dotnetVersion" -ForegroundColor Green
} catch {
    Write-Host ".NET SDK not found. Please install .NET SDK to run integration tests." -ForegroundColor Red
    exit 1
}

$env:INTEGRATION_TESTS = "1"

Write-Host "Running template generation tests..." -ForegroundColor Cyan
go test -v ./internal/generator -run TestTemplateGeneration_

Write-Host "Running folder structure tests..." -ForegroundColor Cyan
go test -v ./internal/generator -run TestFolderStructure_

$env:INTEGRATION_TESTS = $null

Write-Host "Integration tests completed." -ForegroundColor Green 