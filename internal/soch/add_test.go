package soch

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestAddThoughtFromPositionalInput(t *testing.T) {
	notePath := createTestCurrentNote(t, "# Ideas\n\n")
	cmd, output := newTestCommand()
	now := time.Date(2026, time.August, 2, 14, 32, 5, 0, time.FixedZone("EDT", -4*60*60))

	if err := addThoughtAt(cmd, []string{"A useful thought"}, "", now); err != nil {
		t.Fatalf("addThoughtAt() returned an error: %v", err)
	}

	contents, err := os.ReadFile(notePath)
	if err != nil {
		t.Fatal(err)
	}
	want := "# Ideas\n\n- 2026-08-02T14:32:05-04:00 — A useful thought\n"
	if string(contents) != want {
		t.Errorf("note = %q, want %q", contents, want)
	}
	if !strings.Contains(output.String(), "Added thought to "+notePath) {
		t.Error("output did not report the target note")
	}
}

func TestAddThoughtFromFlag(t *testing.T) {
	notePath := createTestCurrentNote(t, "# Ideas\n\n")
	cmd, _ := newTestCommand()
	now := time.Date(2026, time.August, 2, 14, 32, 5, 0, time.UTC)

	if err := addThoughtAt(cmd, nil, "Flag thought", now); err != nil {
		t.Fatalf("addThoughtAt() returned an error: %v", err)
	}

	contents, err := os.ReadFile(notePath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(contents), "2026-08-02T14:32:05Z — Flag thought") {
		t.Errorf("note did not contain flag thought: %q", contents)
	}
}

func TestAddThoughtFromStdin(t *testing.T) {
	notePath := createTestCurrentNote(t, "# Ideas\n\n")
	cmd, _ := newTestCommand()
	cmd.SetIn(strings.NewReader("Piped thought\n"))

	if err := addThoughtAt(cmd, nil, "", time.Now()); err != nil {
		t.Fatalf("addThoughtAt() returned an error: %v", err)
	}

	contents, err := os.ReadFile(notePath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(contents), "— Piped thought\n") {
		t.Errorf("note did not contain stdin thought: %q", contents)
	}
}

func TestAddThoughtPreservesMultilineInput(t *testing.T) {
	notePath := createTestCurrentNote(t, "# Ideas\n\n")
	cmd, _ := newTestCommand()

	if err := addThoughtAt(cmd, []string{"First line\nSecond line"}, "", time.Now()); err != nil {
		t.Fatalf("addThoughtAt() returned an error: %v", err)
	}

	contents, err := os.ReadFile(notePath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(contents), "— First line\n  Second line\n") {
		t.Errorf("note did not preserve multiline thought: %q", contents)
	}
}

func TestAddThoughtRejectsConflictingInputs(t *testing.T) {
	createTestCurrentNote(t, "# Ideas\n\n")
	cmd, _ := newTestCommand()

	err := addThoughtAt(cmd, []string{"Positional"}, "Flag", time.Now())
	if err == nil {
		t.Fatal("addThoughtAt() expected an input conflict")
	}
	if !strings.Contains(err.Error(), "not both") {
		t.Errorf("error = %q, want input guidance", err)
	}
}

func TestAddThoughtRejectsEmptyInput(t *testing.T) {
	createTestCurrentNote(t, "# Ideas\n\n")
	cmd, _ := newTestCommand()
	cmd.SetIn(strings.NewReader(""))

	err := addThoughtAt(cmd, nil, "   ", time.Now())
	if err == nil {
		t.Fatal("addThoughtAt() expected an empty-input error")
	}
	if !strings.Contains(err.Error(), "cannot be empty") {
		t.Errorf("error = %q, want empty-input guidance", err)
	}
}

func TestAddThoughtRequiresCurrentNote(t *testing.T) {
	resetViper(t)
	useTempWorkspace(t)
	if err := os.MkdirAll(projectMetadataDirName, 0755); err != nil {
		t.Fatal(err)
	}

	cmd, _ := newTestCommand()
	err := addThoughtAt(cmd, []string{"Thought"}, "", time.Now())
	if err == nil {
		t.Fatal("addThoughtAt() expected a missing-current-note error")
	}
	if !strings.Contains(err.Error(), "run `soch new") {
		t.Errorf("error = %q, want new-command guidance", err)
	}
}

func TestAddThoughtRejectsBrokenCurrentNote(t *testing.T) {
	resetViper(t)
	useTempWorkspace(t)
	if err := os.MkdirAll(projectMetadataDirName, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join("..", defaultNotesDirName, "missing.md"), currentNotePath); err != nil {
		t.Fatal(err)
	}

	cmd, _ := newTestCommand()
	err := addThoughtAt(cmd, []string{"Thought"}, "", time.Now())
	if err == nil {
		t.Fatal("addThoughtAt() expected a broken-pointer error")
	}
	if !strings.Contains(err.Error(), "does not exist") {
		t.Errorf("error = %q, want broken-pointer guidance", err)
	}
}

func createTestCurrentNote(t *testing.T, contents string) string {
	t.Helper()

	resetViper(t)
	useTempWorkspace(t)
	if err := os.MkdirAll(defaultNotesDirName, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(projectMetadataDirName, 0755); err != nil {
		t.Fatal(err)
	}

	notePath := filepath.Join(defaultNotesDirName, "2026-08-02-ideas.md")
	if err := os.WriteFile(notePath, []byte(contents), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join("..", notePath), currentNotePath); err != nil {
		t.Fatal(err)
	}

	return notePath
}
