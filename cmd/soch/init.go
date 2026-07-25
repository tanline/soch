package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultNotesDirName       = "notes"
	notesDirNameKey           = "notes_dir_name"
	addNotesDirToGitignoreKey = "add_notes_dir_to_gitignore"
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

	if shouldAddNotesDirToGitignore() {
		added, err := addNotesDirToGitignore(notesDirName)
		if err != nil {
			return fmt.Errorf("updating .gitignore: %w", err)
		}

		if added {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Adding %s directory to .gitignore\n", notesDirName); err != nil {
				return err
			}
		}
	}

	_, err := fmt.Fprintln(cmd.OutOrStdout(), "Notes directory is ready")
	return err
}

func shouldAddNotesDirToGitignore() bool {
	if !viper.IsSet(addNotesDirToGitignoreKey) {
		return true
	}

	return viper.GetBool(addNotesDirToGitignoreKey)
}

func addNotesDirToGitignore(notesDirName string) (bool, error) {
	const gitignoreFileName = ".gitignore"

	contents, err := os.ReadFile(gitignoreFileName)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	gitignoreEntry := notesDirName + "/"
	targetEntry := strings.Trim(notesDirName, "/")
	for line := range strings.SplitSeq(strings.ReplaceAll(string(contents), "\r\n", "\n"), "\n") {
		entry := strings.Trim(strings.TrimSpace(line), "/")
		if entry == targetEntry {
			return false, nil
		}
	}

	updatedContents := string(contents)
	if len(updatedContents) > 0 && !strings.HasSuffix(updatedContents, "\n") {
		updatedContents += "\n"
	}
	updatedContents += gitignoreEntry + "\n"

	if err := os.WriteFile(gitignoreFileName, []byte(updatedContents), 0644); err != nil {
		return false, err
	}

	return true, nil
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
