package soch

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:           "soch",
	Short:         "A lightweight CLI that scaffolds, manages, and centralizes notes across multiple repositories without duplicating files.",
	SilenceErrors: true,
	SilenceUsage:  true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return loadConfig(cmd)
	},
}

// Execute adds all child commands and runs the root command.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String(configFlagName, "", "Path to config file")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
