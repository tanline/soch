package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/spf13/cobra"
)

const currentNotePath = ".soch/current-note"

var newCmd = &cobra.Command{
	Use:   "new <title>",
	Short: "Create or select a note collection",
	Args:  cobra.ExactArgs(1),
	RunE:  createNewNote,
}

func createNewNote(cmd *cobra.Command, args []string) error {
	return createNewNoteAt(cmd, args, time.Now())
}

func createNewNoteAt(cmd *cobra.Command, args []string, now time.Time) error {
	notesDirName := getNotesDirName()
	if err := validateNoteWorkspace(notesDirName); err != nil {
		return err
	}
	if err := validateCurrentNotePath(); err != nil {
		return err
	}

	notePath, created, err := createCollection(notesDirName, args[0], now)
	if err != nil {
		return err
	}

	if err := setCurrentNote(notePath); err != nil {
		return fmt.Errorf("selecting current note: %w", err)
	}

	if created {
		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Created collection: %s\n", notePath); err != nil {
			return err
		}
	} else if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Using existing collection: %s\n", notePath); err != nil {
		return err
	}

	_, err = fmt.Fprintf(cmd.OutOrStdout(), "Current note: %s\n", notePath)
	return err
}

func validateNoteWorkspace(notesDirName string) error {
	if err := validateNotesDirName(notesDirName); err != nil {
		return err
	}

	for _, path := range []string{notesDirName, projectMetadataDirName} {
		info, err := os.Stat(path)
		if os.IsNotExist(err) {
			return fmt.Errorf("project is not initialized; run `soch init` before using notes")
		}
		if err != nil {
			return fmt.Errorf("checking project path %q: %w", path, err)
		}
		if !info.IsDir() {
			return fmt.Errorf("project path %q is not a directory; run `soch init` after resolving the conflict", path)
		}
	}

	return nil
}

func createCollection(notesDirName, title string, now time.Time) (string, bool, error) {
	slug := slugifyTitle(title)
	if slug == "" {
		return "", false, fmt.Errorf("title %q does not produce a valid filename", title)
	}

	notePath := filepath.Join(notesDirName, now.Format("2006-01-02")+"-"+slug+".md")
	info, err := os.Stat(notePath)
	if err == nil {
		if info.IsDir() {
			return "", false, fmt.Errorf("collection path %q is a directory", notePath)
		}
		return notePath, false, nil
	}
	if !os.IsNotExist(err) {
		return "", false, fmt.Errorf("checking collection path %q: %w", notePath, err)
	}

	contents := fmt.Sprintf("# %s\n\n", strings.TrimSpace(title))
	if err := os.WriteFile(notePath, []byte(contents), 0644); err != nil {
		return "", false, fmt.Errorf("creating collection %q: %w", notePath, err)
	}

	return notePath, true, nil
}

func slugifyTitle(title string) string {
	var slug strings.Builder
	lastWasSeparator := true

	for _, char := range strings.ToLower(strings.TrimSpace(title)) {
		switch {
		case unicode.IsLetter(char) || unicode.IsDigit(char):
			slug.WriteRune(char)
			lastWasSeparator = false
		case !lastWasSeparator:
			slug.WriteByte('-')
			lastWasSeparator = true
		}
	}

	return strings.Trim(slug.String(), "-")
}

func setCurrentNote(notePath string) error {
	if err := validateCurrentNotePath(); err != nil {
		return err
	}

	info, err := os.Lstat(currentNotePath)
	if err == nil {
		if info.Mode()&os.ModeSymlink == 0 {
			return fmt.Errorf("%q already exists and is not a symlink; remove or rename it", currentNotePath)
		}
		if err := os.Remove(currentNotePath); err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	target, err := filepath.Rel(filepath.Dir(currentNotePath), notePath)
	if err != nil {
		return err
	}

	return os.Symlink(target, currentNotePath)
}

func validateCurrentNotePath() error {
	info, err := os.Lstat(currentNotePath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("checking current note path %q: %w", currentNotePath, err)
	}
	if info.Mode()&os.ModeSymlink == 0 {
		return fmt.Errorf("%q already exists and is not a symlink; remove or rename it", currentNotePath)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(newCmd)
}
