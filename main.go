package main

import (
	"fmt"
	"os"

	"github.com/Andrei-hub11/archforge/cmd"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "archforge",
		Short: "Generator of templates to C#",
		Long:  `Archforge is a tool to generate project structure.`,
	}

	rootCmd.AddCommand(cmd.CreateCmd())
	rootCmd.AddCommand(cmd.InteractiveCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
