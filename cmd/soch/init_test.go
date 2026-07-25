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
	configureMasterNotesDir(t)

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
	configureMasterNotesDir(t)

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

	var output bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&output)

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
	configureMasterNotesDir(t)

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

	var output bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&output)

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
	configureMasterNotesDir(t)

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
	configureMasterNotesDir(t)

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

func TestCreateNotesDirLinksToMasterNotes(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	masterNotesDir := configureMasterNotesDir(t)

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getting current directory: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Errorf("restoring current directory: %v", err)
		}
	})

	workspaceDir := t.TempDir()
	if err := os.Chdir(workspaceDir); err != nil {
		t.Fatalf("changing to temporary directory: %v", err)
	}

	var output bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&output)

	if err := createNotesDir(cmd); err != nil {
		t.Fatalf("createNotesDir() returned an error: %v", err)
	}

	if err := createNotesDir(cmd); err != nil {
		t.Fatalf("second createNotesDir() returned an error: %v", err)
	}

	linkPath := filepath.Join(masterNotesDir, filepath.Base(workspaceDir))
	info, err := os.Lstat(linkPath)
	if err != nil {
		t.Fatalf("checking notes symlink: %v", err)
	}

	if info.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("path %q is not a symlink", linkPath)
	}

	target, err := os.Readlink(linkPath)
	if err != nil {
		t.Fatalf("reading notes symlink: %v", err)
	}

	resolvedTarget, err := filepath.EvalSymlinks(target)
	if err != nil {
		t.Fatalf("resolving symlink target: %v", err)
	}

	wantTarget, err := filepath.EvalSymlinks(filepath.Join(workspaceDir, defaultNotesDirName))
	if err != nil {
		t.Fatalf("resolving expected symlink target: %v", err)
	}

	if resolvedTarget != wantTarget {
		t.Errorf("symlink target = %q, want %q", resolvedTarget, wantTarget)
	}

	if !strings.Contains(output.String(), "Linked notes directory to master notes") {
		t.Error("output did not report linking notes directory to master notes")
	}

	if strings.Count(output.String(), "Linked notes directory to master notes") != 1 {
		t.Error("output reported linking the same notes directory more than once")
	}
}

func configureMasterNotesDir(t *testing.T) string {
	t.Helper()

	masterNotesDir := t.TempDir()
	viper.Set(masterNotesDirKey, masterNotesDir)

	return masterNotesDir
}
