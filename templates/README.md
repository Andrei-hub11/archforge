# Template Files

This directory contains the template files used by the project generator. The templates are organized in subdirectories based on the project type:

- `console` - Templates for console applications
- `webapi` - Templates for Web API applications
- `mvc` - Templates for MVC applications

## Template Format

Template files support Go template syntax for dynamic content replacement. The following variables are available in templates:

- `{{.ProjectName}}` - The name of the project being generated

You can use these variables both in file content and filenames. For filenames, any occurrence of `{{.ProjectName}}` will be replaced with the actual project name during generation.

## Adding New Templates

To add new templates:

1. Create a new subdirectory for your template type
2. Add all the necessary files with appropriate template variables
3. Add a new generator function in `internal/templates/template.go`
4. Update the `Generate` function in `internal/generator/generator.go` to call your new generator

## Example

File: `templates/console/{{.ProjectName}}.csproj`

```xml
<Project Sdk="Microsoft.NET.Sdk">
  <PropertyGroup>
    <OutputType>Exe</OutputType>
    <TargetFramework>net7.0</TargetFramework>
    <ImplicitUsings>enable</ImplicitUsings>
    <Nullable>enable</Nullable>
  </PropertyGroup>
</Project>
```

This will create a .csproj file with the project name when generating a console application.
