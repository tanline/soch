package soch

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show [filename]",
	Short: "Show the current note or a named note",
	Args:  cobra.MaximumNArgs(1),
	RunE:  showNote,
}

func showNote(cmd *cobra.Command, args []string) error {
	notePath, err := getNotePathToShow(args)
	if err != nil {
		return err
	}

	contents, err := os.ReadFile(notePath)
	if err != nil {
		return fmt.Errorf("reading note %q: %w", notePath, err)
	}

	_, err = cmd.OutOrStdout().Write(contents)
	return err
}

func getNotePathToShow(args []string) (string, error) {
	if len(args) == 0 {
		return getCurrentNotePath()
	}

	return getNotePathByName(args[0])
}

func getNotePathByName(filename string) (string, error) {
	if err := validateNoteFilename(filename); err != nil {
		return "", err
	}

	notesDirName := getNotesDirName()
	if err := validateNotesDirName(notesDirName); err != nil {
		return "", err
	}

	notePath := filepath.Join(notesDirName, filename)
	info, err := os.Lstat(notePath)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("note %q does not exist; create it with `soch new \"<title>\"` first", filename)
	}
	if err != nil {
		return "", fmt.Errorf("checking note %q: %w", filename, err)
	}
	if !info.Mode().IsRegular() {
		return "", fmt.Errorf("note %q is not a regular file", filename)
	}

	return notePath, nil
}

func validateNoteFilename(filename string) error {
	if filename == "" || filename == "." || filename == ".." ||
		filepath.IsAbs(filename) || filepath.Base(filename) != filename {
		return fmt.Errorf("invalid note filename %q; use a single filename such as `ideas.md`", filename)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(showCmd)
}
