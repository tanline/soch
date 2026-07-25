package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultNotesDirName       = "notes"
	notesDirNameKey           = "notes_dir_name"
	addNotesDirToGitignoreKey = "add_notes_dir_to_gitignore"
	masterNotesDirKey         = "master_notes_dir"
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

	if err := createLocalNotesDir(notesDirName); err != nil {
		return fmt.Errorf("creating notes directory: %w", err)
	}

	if shouldAddNotesDirToGitignore() {
		added, err := updateGitignore(notesDirName)
		if err != nil {
			return fmt.Errorf("updating .gitignore: %w", err)
		}

		if added {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Adding %s directory to .gitignore\n", notesDirName); err != nil {
				return err
			}
		}
	}

	linked, err := linkNotesDir(notesDirName)
	if err != nil {
		return fmt.Errorf("linking notes directory: %w", err)
	}

	if linked {
		if _, err := fmt.Fprintln(cmd.OutOrStdout(), "Linked notes directory to master notes"); err != nil {
			return err
		}
	}

	_, err = fmt.Fprintln(cmd.OutOrStdout(), "Notes directory is ready")
	return err
}

func createLocalNotesDir(notesDirName string) error {
	return os.MkdirAll(notesDirName, 0755)
}

func linkNotesDir(notesDirName string) (bool, error) {
	masterNotesDir := strings.TrimSpace(viper.GetString(masterNotesDirKey))
	if masterNotesDir == "" {
		return false, fmt.Errorf("master notes directory is not configured")
	}

	info, err := os.Stat(masterNotesDir)
	if err != nil {
		return false, fmt.Errorf("accessing master notes directory %q: %w", masterNotesDir, err)
	}
	if !info.IsDir() {
		return false, fmt.Errorf("master notes path %q is not a directory", masterNotesDir)
	}

	workspaceDir, err := os.Getwd()
	if err != nil {
		return false, fmt.Errorf("getting current workspace: %w", err)
	}

	linkPath := filepath.Join(masterNotesDir, filepath.Base(workspaceDir))
	notesPath := filepath.Join(workspaceDir, notesDirName)

	linkInfo, err := os.Lstat(linkPath)
	if err == nil {
		if linkInfo.Mode()&os.ModeSymlink == 0 {
			return false, fmt.Errorf("link path %q already exists and is not a symlink", linkPath)
		}

		existingTarget, err := os.Readlink(linkPath)
		if err != nil {
			return false, fmt.Errorf("reading existing symlink: %w", err)
		}

		resolvedExistingTarget, err := filepath.EvalSymlinks(existingTarget)
		if err != nil {
			return false, fmt.Errorf("resolving existing symlink target: %w", err)
		}

		resolvedNotesPath, err := filepath.EvalSymlinks(notesPath)
		if err != nil {
			return false, fmt.Errorf("resolving notes directory: %w", err)
		}

		if resolvedExistingTarget == resolvedNotesPath {
			return false, nil
		}

		return false, fmt.Errorf("link path %q points to a different directory", linkPath)
	}
	if !os.IsNotExist(err) {
		return false, fmt.Errorf("checking link path: %w", err)
	}

	if err := os.Symlink(notesPath, linkPath); err != nil {
		return false, err
	}

	return true, nil
}

func shouldAddNotesDirToGitignore() bool {
	if !viper.IsSet(addNotesDirToGitignoreKey) {
		return true
	}

	return viper.GetBool(addNotesDirToGitignoreKey)
}

func updateGitignore(notesDirName string) (bool, error) {
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
	}

	return name
}

func init() {
	rootCmd.AddCommand(initCmd)
}
