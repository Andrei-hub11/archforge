package cmd

import (
	"fmt"
	"os"
	"strings"

	"errors"

	"github.com/Andrei-hub11/archforge/internal/config"
	"github.com/Andrei-hub11/archforge/internal/generator"
	"github.com/spf13/cobra"
)

func CreateCmd() *cobra.Command {
	var cfg config.ProjectConfig

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new project",
		Run: func(cmd *cobra.Command, args []string) {
			var missing []string

			if cfg.Name == "" {
				missing = append(missing, "--name")
			}

			if cfg.Template == "" {
				missing = append(missing, "--template")
			}

			if len(missing) > 0 {
				fmt.Printf("Erro: faltando flag(s) obrigatória(s): %s\n\n", strings.Join(missing, ", "))
				cmd.Usage()
				os.Exit(1)
			}

			if err := validateProjectName(cfg.Name); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if cfg.Preview {
				fmt.Println("Previewing the project...")
				err := generator.GenerateBuildTree(cfg)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				if getResponseToPreview() {
					err := generator.Generate(cfg)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
				}

				os.Exit(0)
			} else {
				err := generator.Generate(cfg)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
		},
	}

	createCmd.Flags().StringVarP(&cfg.Name, "name", "n", "", "The name of the project")
	createCmd.Flags().StringVarP(&cfg.OutputDir, "output", "o", "", "The output directory (optional, defaults to path)")
	createCmd.Flags().StringVarP(&cfg.Template, "template", "t", "", "The template to use")
	createCmd.Flags().BoolVarP(&cfg.Preview, "preview", "p", false, "Preview the project without generating it")

	createCmd.MarkFlagRequired("name")
	createCmd.MarkFlagRequired("template")

	return createCmd
}

func validateProjectName(name string) error {
	if name == "" {
		return errors.New("project name is required")
	}

	if strings.Contains(name, " ") {
		return errors.New("project name must not contain spaces")
	}

	if strings.Contains(name, "-") {
		return errors.New("project name must not contain hyphens")
	}

	return nil
}
