package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultNotesDirName = "notes"
	notesDirNameKey     = "notes_dir_name"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a notes directory in current workspace",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := createNotesDir(cmd)
		return err
	},
}

func createNotesDir(cmd *cobra.Command) error {
	notesDirName := getNotesDirName()

	if err := os.MkdirAll(notesDirName, 0755); err != nil {
		return fmt.Errorf("creating notes directory: %w", err)
	}

	_, err := fmt.Fprintln(cmd.OutOrStdout(), "Notes directory is ready")
	return err
}

func getNotesDirName() string {
	name := strings.TrimSpace(viper.GetString(notesDirNameKey))

	if name == "" {
		return defaultNotesDirName
	} else {
		return name
	}
}

func init() {
	rootCmd.AddCommand(initCmd)
}
