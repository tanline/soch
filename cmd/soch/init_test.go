package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestCreateNotesDir(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)

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

	if err := os.WriteFile(".gitignore", []byte("*.log\n"), 0644); err != nil {
		t.Fatalf("creating .gitignore: %v", err)
	}

	var output bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&output)

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

	contents, err := os.ReadFile(".gitignore")
	if err != nil {
		t.Fatalf("reading .gitignore: %v", err)
	}

	wantGitignore := "*.log\nnotes/\n"
	if string(contents) != wantGitignore {
		t.Errorf(".gitignore = %q, want %q", string(contents), wantGitignore)
	}

	if !strings.Contains(output.String(), "Adding notes directory to .gitignore") {
		t.Error("output did not report adding notes directory to .gitignore")
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

func TestCreateNotesDirDoesNotCreateGitignore(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getting current directory: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Errorf("restoring current directory: %v", err)
		}
	})

	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatalf("changing to temporary directory: %v", err)
	}

	cmd := &cobra.Command{}
	cmd.SetOut(&bytes.Buffer{})

	if err := createNotesDir(cmd); err != nil {
		t.Fatalf("createNotesDir() returned an error: %v", err)
	}

	if _, err := os.Stat(".gitignore"); !os.IsNotExist(err) {
		t.Fatalf(".gitignore was created or could not be checked: %v", err)
	}
}

func TestCreateNotesDirDoesNotDuplicateGitignoreEntry(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getting current directory: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Errorf("restoring current directory: %v", err)
		}
	})

	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatalf("changing to temporary directory: %v", err)
	}

	if err := os.WriteFile(".gitignore", []byte("notes\n"), 0644); err != nil {
		t.Fatalf("creating .gitignore: %v", err)
	}

	cmd := &cobra.Command{}
	cmd.SetOut(&bytes.Buffer{})

	if err := createNotesDir(cmd); err != nil {
		t.Fatalf("createNotesDir() returned an error: %v", err)
	}

	contents, err := os.ReadFile(".gitignore")
	if err != nil {
		t.Fatalf("reading .gitignore: %v", err)
	}

	if string(contents) != "notes\n" {
		t.Errorf(".gitignore = %q, want %q", string(contents), "notes\n")
	}
}

func TestCreateNotesDirSkipsGitignoreWhenDisabled(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.Set(addNotesDirToGitignoreKey, false)

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getting current directory: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Errorf("restoring current directory: %v", err)
		}
	})

	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatalf("changing to temporary directory: %v", err)
	}

	if err := os.WriteFile(".gitignore", []byte("*.log\n"), 0644); err != nil {
		t.Fatalf("creating .gitignore: %v", err)
	}

	var output bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&output)

	if err := createNotesDir(cmd); err != nil {
		t.Fatalf("createNotesDir() returned an error: %v", err)
	}

	contents, err := os.ReadFile(".gitignore")
	if err != nil {
		t.Fatalf("reading .gitignore: %v", err)
	}

	if string(contents) != "*.log\n" {
		t.Errorf(".gitignore = %q, want %q", string(contents), "*.log\n")
	}

	if strings.Contains(output.String(), "Adding notes directory to .gitignore") {
		t.Error("output reported adding notes directory to .gitignore")
	}
}
