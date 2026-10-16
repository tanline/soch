package soch

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newSetupCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "setup",
		Short: "Configure the global settings for soch",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := initSochConfig(cmd); err != nil {
				return fmt.Errorf("initializing soch config: %w", err)
			}

			if err := setMasterNotesDir(cmd); err != nil {
				return fmt.Errorf("setting master notes dir: %w", err)
			}

			if err := viper.WriteConfig(); err != nil {
				return fmt.Errorf("writing soch config: %w", err)
			}

			_, err := fmt.Fprintln(cmd.OutOrStdout(), "Finished setting up soch!")
			return err
		},
	}
}

func setMasterNotesDir(cmd *cobra.Command) error {
	_, err := fmt.Fprint(cmd.OutOrStdout(), "Enter the absolute path to your master notes directory: ")
	if err != nil {
		return err
	}

	reader := bufio.NewReader(cmd.InOrStdin())
	inputPath, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	inputPath = strings.TrimSpace(inputPath)

	if strings.HasPrefix(inputPath, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("finding user home directory: %w", err)
		}

		inputPath = filepath.Join(homeDir, inputPath[2:])
	}

	if !filepath.IsAbs(inputPath) {
		return fmt.Errorf("master notes path must be absolute")
	}

	info, err := os.Stat(inputPath)
	if err != nil {
		return fmt.Errorf("cannot access directory %q: %w", inputPath, err)
	}

	if !info.IsDir() {
		return fmt.Errorf("path %q is not a directory", inputPath)
	}

	viper.Set("master_notes_dir", inputPath)

	_, err = fmt.Fprintln(cmd.OutOrStdout(), "Your master notes directory is located at:", inputPath)
	return err
}
