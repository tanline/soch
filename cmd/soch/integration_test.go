package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLIWorkflow(t *testing.T) {
	binaryPath := buildTestBinary(t)
	temporaryDir := t.TempDir()
	projectDir := filepath.Join(temporaryDir, "project")
	masterNotesDir := filepath.Join(temporaryDir, "master-notes")
	configPath := filepath.Join(temporaryDir, "config.yaml")

	for _, directory := range []string{projectDir, masterNotesDir} {
		if err := os.MkdirAll(directory, 0755); err != nil {
			t.Fatalf("creating test directory %q: %v", directory, err)
		}
	}
	writeIntegrationConfig(t, configPath, masterNotesDir)

	runSoch(t, binaryPath, projectDir, configPath, "init")
	runSoch(t, binaryPath, projectDir, configPath, "new", "Ideas")
	runSoch(t, binaryPath, projectDir, configPath, "add", "A test thought")

	notePath := findCreatedNote(t, filepath.Join(projectDir, "project-notes"))
	noteContents, err := os.ReadFile(notePath)
	if err != nil {
		t.Fatalf("reading created note: %v", err)
	}
	if !strings.Contains(string(noteContents), "# Ideas\n") {
		t.Errorf("note contents = %q, want title", noteContents)
	}
	if !strings.Contains(string(noteContents), "A test thought") {
		t.Errorf("note contents = %q, want added thought", noteContents)
	}

	showOutput := runSoch(t, binaryPath, projectDir, configPath, "show")
	if showOutput != string(noteContents) {
		t.Errorf("show output = %q, want %q", showOutput, noteContents)
	}

	assertCurrentNoteTarget(t, projectDir, notePath)
	assertMasterNotesLink(t, masterNotesDir, projectDir, "project-notes")
}

func buildTestBinary(t *testing.T) string {
	t.Helper()

	binaryPath := filepath.Join(t.TempDir(), "soch")
	command := exec.Command("go", "build", "-o", binaryPath, ".")
	command.Dir = "."
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("building test binary: %v\n%s", err, output)
	}

	return binaryPath
}

func writeIntegrationConfig(t *testing.T, configPath, masterNotesDir string) {
	t.Helper()

	templatePath := filepath.Join("testdata", "integration-config.yaml.tmpl")
	contents, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("reading integration config fixture: %v", err)
	}

	config := strings.ReplaceAll(string(contents), "__MASTER_NOTES_DIR__", masterNotesDir)
	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		t.Fatalf("writing integration config: %v", err)
	}
}

func runSoch(t *testing.T, binaryPath, projectDir, configPath string, args ...string) string {
	t.Helper()

	commandArgs := append([]string{"--config", configPath}, args...)
	command := exec.Command(binaryPath, commandArgs...)
	command.Dir = projectDir
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("running soch %s: %v\n%s", strings.Join(args, " "), err, output)
	}

	return string(output)
}

func findCreatedNote(t *testing.T, notesDir string) string {
	t.Helper()

	entries, err := os.ReadDir(notesDir)
	if err != nil {
		t.Fatalf("reading notes directory: %v", err)
	}

	var noteNames []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			noteNames = append(noteNames, entry.Name())
		}
	}
	if len(noteNames) != 1 {
		t.Fatalf("notes directory contains %d markdown files, want 1", len(noteNames))
	}

	return filepath.Join(notesDir, noteNames[0])
}

func assertCurrentNoteTarget(t *testing.T, projectDir, notePath string) {
	t.Helper()

	currentNote := filepath.Join(projectDir, ".soch", "current-note")
	resolvedCurrentNote, err := filepath.EvalSymlinks(currentNote)
	if err != nil {
		t.Fatalf("resolving current note symlink: %v", err)
	}
	resolvedNotePath, err := filepath.EvalSymlinks(notePath)
	if err != nil {
		t.Fatalf("resolving note path: %v", err)
	}
	if resolvedCurrentNote != resolvedNotePath {
		t.Errorf("current note target = %q, want %q", resolvedCurrentNote, resolvedNotePath)
	}
}

func assertMasterNotesLink(t *testing.T, masterNotesDir, projectDir, notesDirName string) {
	t.Helper()

	linkPath := filepath.Join(masterNotesDir, filepath.Base(projectDir))
	resolvedLink, err := filepath.EvalSymlinks(linkPath)
	if err != nil {
		t.Fatalf("resolving master notes link: %v", err)
	}
	resolvedNotesDir, err := filepath.EvalSymlinks(filepath.Join(projectDir, notesDirName))
	if err != nil {
		t.Fatalf("resolving notes directory: %v", err)
	}
	if resolvedLink != resolvedNotesDir {
		t.Errorf("master notes link target = %q, want %q", resolvedLink, resolvedNotesDir)
	}
}
