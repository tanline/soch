package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestCreateNotesDir(t *testing.T) {
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getting current directory: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Errorf("restoring current directory: %v", err)
		}
	})

	temporaryDir := t.TempDir()
	if err := os.Chdir(temporaryDir); err != nil {
		t.Fatalf("changing to temporary directory: %v", err)
	}

	cmd := &cobra.Command{}
	cmd.SetOut(&bytes.Buffer{})

	if err := createNotesDir(cmd); err != nil {
		t.Fatalf("createNotesDir() returned an error: %v", err)
	}

	info, err := os.Stat(filepath.Join(temporaryDir, defaultNotesDirName))
	if err != nil {
		t.Fatalf("checking notes directory: %v", err)
	}

	if !info.IsDir() {
		t.Fatal("notes path exists but is not a directory")
	}
}

func TestCreateNotesDirUsesConfiguredName(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)

	if err := loadConfigFrom("testdata"); err != nil {
		t.Fatalf("loadConfigFrom() returned an error: %v", err)
	}

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getting current directory: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Errorf("restoring current directory: %v", err)
		}
	})

	temporaryDir := t.TempDir()
	if err := os.Chdir(temporaryDir); err != nil {
		t.Fatalf("changing to temporary directory: %v", err)
	}

	cmd := &cobra.Command{}
	cmd.SetOut(&bytes.Buffer{})

	if err := createNotesDir(cmd); err != nil {
		t.Fatalf("createNotesDir() returned an error: %v", err)
	}

	info, err := os.Stat(filepath.Join(temporaryDir, "project-notes"))
	if err != nil {
		t.Fatalf("checking configured notes directory: %v", err)
	}

	if !info.IsDir() {
		t.Fatal("configured notes path exists but is not a directory")
	}
}
