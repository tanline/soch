package main

import (
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestLoadConfigFrom(t *testing.T) {
	resetViper(t)

	configDir := filepath.Join("testdata")

	if err := loadConfigFrom(configDir); err != nil {
		t.Fatalf("loadConfigFrom() returned an error: %v", err)
	}

	if got := viper.GetString(notesDirNameKey); got != "project-notes" {
		t.Errorf("notes_dir_name = %q, want %q", got, "project-notes")
	}
}
