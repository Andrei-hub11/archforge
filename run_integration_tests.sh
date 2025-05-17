#!/bin/bash
# Bash script to run integration tests for Archforge
# This script will set the INTEGRATION_TESTS environment variable and run the Go tests

echo -e "\033[0;32mRunning Archforge integration tests...\033[0m"

if command -v dotnet &> /dev/null; then
    echo -e "\033[0;32mFound .NET SDK version: $(dotnet --version)\033[0m"
else
    echo -e "\033[0;31m.NET SDK not found. Please install .NET SDK to run integration tests.\033[0m"
    exit 1
fi

export INTEGRATION_TESTS=1

echo -e "\033[0;36mRunning template generation tests...\033[0m"
go test -v ./internal/generator -run TestTemplateGeneration_

echo -e "\033[0;36mRunning folder structure tests...\033[0m"
go test -v ./internal/generator -run TestFolderStructure_

unset INTEGRATION_TESTS

echo -e "\033[0;32mIntegration tests completed.\033[0m" 