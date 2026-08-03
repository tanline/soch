package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestCreateNotesDir(t *testing.T) {
	resetViper(t)
	configureMasterNotesDir(t)
	temporaryDir := useTempWorkspace(t)

	if err := os.WriteFile(".gitignore", []byte("*.log\n"), 0644); err != nil {
		t.Fatalf("creating .gitignore: %v", err)
	}

	cmd, output := newTestCommand()

	if err := createNotesDir(cmd, nil); err != nil {
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
	resetViper(t)

	if err := loadConfigFrom("testdata"); err != nil {
		t.Fatalf("loadConfigFrom() returned an error: %v", err)
	}
	configureMasterNotesDir(t)
	temporaryDir := useTempWorkspace(t)

	cmd, _ := newTestCommand()

	if err := createNotesDir(cmd, nil); err != nil {
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
	resetViper(t)
	configureMasterNotesDir(t)
	useTempWorkspace(t)

	cmd, _ := newTestCommand()

	if err := createNotesDir(cmd, nil); err != nil {
		t.Fatalf("createNotesDir() returned an error: %v", err)
	}

	if _, err := os.Stat(".gitignore"); !os.IsNotExist(err) {
		t.Fatalf(".gitignore was created or could not be checked: %v", err)
	}
}

func TestCreateNotesDirDoesNotDuplicateGitignoreEntry(t *testing.T) {
	resetViper(t)
	configureMasterNotesDir(t)
	useTempWorkspace(t)

	if err := os.WriteFile(".gitignore", []byte("notes\n"), 0644); err != nil {
		t.Fatalf("creating .gitignore: %v", err)
	}

	cmd, _ := newTestCommand()

	if err := createNotesDir(cmd, nil); err != nil {
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
	resetViper(t)
	viper.Set(addNotesDirToGitignoreKey, false)
	configureMasterNotesDir(t)
	useTempWorkspace(t)

	if err := os.WriteFile(".gitignore", []byte("*.log\n"), 0644); err != nil {
		t.Fatalf("creating .gitignore: %v", err)
	}

	cmd, output := newTestCommand()

	if err := createNotesDir(cmd, nil); err != nil {
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
	resetViper(t)
	masterNotesDir := configureMasterNotesDir(t)
	workspaceDir := useTempWorkspace(t)

	cmd, output := newTestCommand()

	if err := createNotesDir(cmd, nil); err != nil {
		t.Fatalf("createNotesDir() returned an error: %v", err)
	}

	if err := createNotesDir(cmd, nil); err != nil {
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

func TestCreateNotesDirRejectsUnsafeName(t *testing.T) {
	tests := []string{"../notes", "/tmp/notes", ".", ".."}

	for _, notesDirName := range tests {
		t.Run(notesDirName, func(t *testing.T) {
			resetViper(t)
			viper.Set(notesDirNameKey, notesDirName)
			configureMasterNotesDir(t)
			useTempWorkspace(t)

			if err := createNotesDir(&cobra.Command{}, nil); err == nil {
				t.Fatal("createNotesDir() expected an error for an unsafe directory name")
			}
		})
	}
}

func TestCreateNotesDirReportsMissingMasterNotesDirBeforeMutation(t *testing.T) {
	resetViper(t)
	workspaceDir := useTempWorkspace(t)

	if err := os.WriteFile(".gitignore", []byte("*.log\n"), 0644); err != nil {
		t.Fatalf("creating .gitignore: %v", err)
	}

	err := createNotesDir(&cobra.Command{}, nil)
	if err == nil {
		t.Fatal("createNotesDir() expected an error")
	}
	if !strings.Contains(err.Error(), "run `soch setup` first") {
		t.Errorf("error = %q, want setup guidance", err)
	}

	if _, err := os.Stat(filepath.Join(workspaceDir, defaultNotesDirName)); !os.IsNotExist(err) {
		t.Errorf("notes directory was created after preflight failure: %v", err)
	}

	contents, err := os.ReadFile(".gitignore")
	if err != nil {
		t.Fatalf("reading .gitignore: %v", err)
	}
	if string(contents) != "*.log\n" {
		t.Errorf(".gitignore = %q, want %q", string(contents), "*.log\n")
	}
}

func TestCreateNotesDirReportsConflictingLinkBeforeMutation(t *testing.T) {
	resetViper(t)
	masterNotesDir := configureMasterNotesDir(t)
	workspaceDir := useTempWorkspace(t)

	if err := os.WriteFile(".gitignore", []byte("*.log\n"), 0644); err != nil {
		t.Fatalf("creating .gitignore: %v", err)
	}

	linkPath := filepath.Join(masterNotesDir, filepath.Base(workspaceDir))
	if err := os.WriteFile(linkPath, []byte("existing"), 0644); err != nil {
		t.Fatalf("creating conflicting link path: %v", err)
	}

	err := createNotesDir(&cobra.Command{}, nil)
	if err == nil {
		t.Fatal("createNotesDir() expected an error")
	}
	if !strings.Contains(err.Error(), "remove or rename it") {
		t.Errorf("error = %q, want remediation guidance", err)
	}

	if _, err := os.Stat(defaultNotesDirName); !os.IsNotExist(err) {
		t.Errorf("notes directory was created after preflight failure: %v", err)
	}

	contents, err := os.ReadFile(".gitignore")
	if err != nil {
		t.Fatalf("reading .gitignore: %v", err)
	}
	if string(contents) != "*.log\n" {
		t.Errorf(".gitignore = %q, want %q", string(contents), "*.log\n")
	}
}

func TestValidateExistingLink(t *testing.T) {
	temporaryDir := t.TempDir()
	expectedTarget := filepath.Join(temporaryDir, "notes")
	linkPath := filepath.Join(temporaryDir, "project")

	exists, err := validateExistingLink(linkPath, expectedTarget)
	if err != nil {
		t.Fatalf("validateExistingLink() returned an error for a missing link: %v", err)
	}
	if exists {
		t.Fatal("validateExistingLink() reported a missing link as existing")
	}

	if err := os.Mkdir(expectedTarget, 0755); err != nil {
		t.Fatalf("creating expected target: %v", err)
	}
	if err := os.Symlink(expectedTarget, linkPath); err != nil {
		t.Fatalf("creating expected symlink: %v", err)
	}

	exists, err = validateExistingLink(linkPath, expectedTarget)
	if err != nil {
		t.Fatalf("validateExistingLink() returned an error for the correct link: %v", err)
	}
	if !exists {
		t.Fatal("validateExistingLink() did not report the correct link as existing")
	}

	if err := os.Remove(linkPath); err != nil {
		t.Fatalf("removing expected symlink: %v", err)
	}
	if err := os.Symlink(filepath.Join(temporaryDir, "other"), linkPath); err != nil {
		t.Fatalf("creating conflicting symlink: %v", err)
	}

	if _, err := validateExistingLink(linkPath, expectedTarget); err == nil {
		t.Fatal("validateExistingLink() expected an error for a conflicting link")
	}

	if err := os.Remove(linkPath); err != nil {
		t.Fatalf("removing conflicting symlink: %v", err)
	}
	if err := os.Symlink("notes", linkPath); err != nil {
		t.Fatalf("creating relative symlink: %v", err)
	}

	exists, err = validateExistingLink(linkPath, expectedTarget)
	if err != nil {
		t.Fatalf("validateExistingLink() returned an error for a relative link: %v", err)
	}
	if !exists {
		t.Fatal("validateExistingLink() did not report the relative link as existing")
	}
}

func configureMasterNotesDir(t *testing.T) string {
	t.Helper()

	masterNotesDir := t.TempDir()
	viper.Set(masterNotesDirKey, masterNotesDir)

	return masterNotesDir
}
