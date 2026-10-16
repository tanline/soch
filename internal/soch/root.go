package soch

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "soch",
		Short:         "A lightweight CLI that scaffolds, manages, and centralizes notes across multiple repositories without duplicating files.",
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return loadConfig(cmd)
		},
	}

	rootCmd.PersistentFlags().String(configFlagName, "", "Path to config file")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(newSetupCommand(), newInitCommand(), newNoteCommand(), newAddCommand(), newShowCommand())

	return rootCmd
}

// Execute runs the CLI and reports command errors.
func Execute() {
	err := newRootCommand().Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
