package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Andrei-hub11/archforge/internal/config"
	"github.com/Andrei-hub11/archforge/internal/datas"
	"github.com/Andrei-hub11/archforge/internal/generator"
	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
)

// InputFunc defines the function signature for getting user input
type InputFunc func(message string, defaultValue string) string

// SelectFunc defines the function signature for selecting from options
type SelectFunc func(message string, options []string, defaultValue string) string

// Default implementation using go-prompt
var defaultInputFunc = func(message string, defaultValue string) string {
	fmt.Printf("%s (default: %s): ", message, defaultValue)

	completer := func(d prompt.Document) []prompt.Suggest {
		return []prompt.Suggest{}
	}

	options := []prompt.Option{
		prompt.OptionPrefixTextColor(prompt.Blue),
		prompt.OptionInputTextColor(prompt.DefaultColor),
		prompt.OptionCompletionWordSeparator(" "),
	}

	result := prompt.Input("", completer, options...)
	if result == "" {
		return defaultValue
	}
	return result
}

// Default implementation for selection using go-prompt
var defaultSelectFunc = func(message string, options []string, defaultValue string) string {
	fmt.Printf("%s (default: %s): ", message, defaultValue)

	// Create suggestions for each option
	suggestions := []prompt.Suggest{}
	for _, opt := range options {
		suggestions = append(suggestions, prompt.Suggest{Text: opt})
	}

	// Completer that suggests valid options
	completer := func(d prompt.Document) []prompt.Suggest {
		return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
	}

	promptOptions := []prompt.Option{
		prompt.OptionPrefixTextColor(prompt.Blue),
		prompt.OptionInputTextColor(prompt.DefaultColor),
		prompt.OptionCompletionWordSeparator(" "),
		prompt.OptionShowCompletionAtStart(),
	}

	result := prompt.Input("", completer, promptOptions...)

	// Validate if the result is one of the options
	valid := false
	for _, opt := range options {
		if strings.HasPrefix(opt, result) || strings.HasPrefix(result, strings.Split(opt, " - ")[0]) {
			valid = true
			result = strings.Split(opt, " - ")[0]
			break
		}
	}

	if result == "" || !valid {
		return defaultValue
	}

	return result
}

// Current functions that can be changed for testing
var currentInputFunc InputFunc = defaultInputFunc
var currentSelectFunc SelectFunc = defaultSelectFunc

// SetInputFunc allows changing the input function for testing
func SetInputFunc(inputFunc InputFunc) {
	currentInputFunc = inputFunc
}

// SetSelectFunc allows changing the select function for testing
func SetSelectFunc(selectFunc SelectFunc) {
	currentSelectFunc = selectFunc
}

// ResetInputFunc restores the default input function
func ResetInputFunc() {
	currentInputFunc = defaultInputFunc
}

// ResetSelectFunc restores the default select function
func ResetSelectFunc() {
	currentSelectFunc = defaultSelectFunc
}

// InteractiveCmd returns the interactive command
func InteractiveCmd() *cobra.Command {
	interactiveCmd := &cobra.Command{
		Use:   "interactive",
		Short: "Interactive mode",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := getInteractiveConfig()

			if err := validateProjectName(cfg.Name); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if cfg.Preview {
				fmt.Println("Previewing project...")
				err := generator.GenerateBuildTree(cfg)
				if err != nil {
					fmt.Println(err)
				}

				if getResponseToPreview() {
					err := generator.Generate(cfg)
					if err != nil {
						log.Fatal(err)
					}
				}
			} else {
				err := generator.Generate(cfg)
				if err != nil {
					log.Fatal(err)
				}
			}
		},
	}

	return interactiveCmd
}

func getInteractiveConfig() config.ProjectConfig {
	defaultName := "MyProject"
	defaultTemplate := "console"
	defaultOutput := "."

	cfg := config.ProjectConfig{
		Name:      defaultName,
		Template:  defaultTemplate,
		OutputDir: defaultOutput,
	}

	cfg.Name = currentInputFunc("Enter the name of the project", defaultName)

	output := currentInputFunc("Enter the output directory", defaultOutput)
	if output != "" {
		cfg.OutputDir = output
	}

	templateID := currentSelectFunc(
		"Select the template to use",
		datas.GetTemplateOptions(),
		"1",
	)

	cfg.Template = datas.TemplatesSelect[templateID]

	// Get preview option using the select function
	preview := currentSelectFunc("Do you want to preview the project? (1/Yes, 2/No)", []string{"1", "2"}, "2")
	cfg.Preview = preview == "1"

	return cfg
}
func getResponseToPreview() bool {
	response := currentSelectFunc("Do you accept the preview? (1/Yes, 2/No)", []string{"1", "2"}, "2")
	return response == "1"
}
