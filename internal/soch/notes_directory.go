package soch

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const (
	defaultNotesDirName = "notes"
	notesDirNameKey     = "notes_dir_name"
)

func getNotesDirName() string {
	name := strings.TrimSpace(viper.GetString(notesDirNameKey))
	if name == "" {
		return defaultNotesDirName
	}

	return name
}

func validateNotesDirName(notesDirName string) error {
	if notesDirName == "." || notesDirName == ".." ||
		filepath.IsAbs(notesDirName) || filepath.Base(notesDirName) != notesDirName {
		return fmt.Errorf("invalid notes directory name %q; use a single relative directory name such as `notes`", notesDirName)
	}

	return nil
}
