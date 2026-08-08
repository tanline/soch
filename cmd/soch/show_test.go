package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestShowNoteOutputsCurrentNote(t *testing.T) {
	resetViper(t)
	useTempWorkspace(t)
	if err := os.MkdirAll(projectMetadataDirName, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(defaultNotesDirName, 0755); err != nil {
		t.Fatal(err)
	}

	notePath := filepath.Join(defaultNotesDirName, "ideas.md")
	contents := "# Ideas\n\n- A thought\n"
	if err := os.WriteFile(notePath, []byte(contents), 0644); err != nil {
		t.Fatal(err)
	}
	if err := setCurrentNote(notePath); err != nil {
		t.Fatal(err)
	}

	cmd, output := newTestCommand()
	if err := showNote(cmd, nil); err != nil {
		t.Fatalf("showNote() returned an error: %v", err)
	}

	if output.String() != contents {
		t.Errorf("output = %q, want %q", output.String(), contents)
	}
}

func TestShowNoteOutputsNamedNoteWithoutCurrentNote(t *testing.T) {
	resetViper(t)
	useTempWorkspace(t)
	if err := os.MkdirAll(defaultNotesDirName, 0755); err != nil {
		t.Fatal(err)
	}

	contents := "# Ideas\n\n- A thought\n"
	if err := os.WriteFile(filepath.Join(defaultNotesDirName, "ideas.md"), []byte(contents), 0644); err != nil {
		t.Fatal(err)
	}

	cmd, output := newTestCommand()
	if err := showNote(cmd, []string{"ideas.md"}); err != nil {
		t.Fatalf("showNote() returned an error: %v", err)
	}

	if output.String() != contents {
		t.Errorf("output = %q, want %q", output.String(), contents)
	}
}

func TestShowNoteUsesConfiguredNotesDirectory(t *testing.T) {
	resetViper(t)
	useTempWorkspace(t)
	viper.Set(notesDirNameKey, "project-notes")
	if err := os.MkdirAll("project-notes", 0755); err != nil {
		t.Fatal(err)
	}

	contents := "# Ideas\n\n- A thought\n"
	if err := os.WriteFile(filepath.Join("project-notes", "ideas.md"), []byte(contents), 0644); err != nil {
		t.Fatal(err)
	}

	cmd, output := newTestCommand()
	if err := showNote(cmd, []string{"ideas.md"}); err != nil {
		t.Fatalf("showNote() returned an error: %v", err)
	}

	if output.String() != contents {
		t.Errorf("output = %q, want %q", output.String(), contents)
	}
}

func TestShowNoteRequiresCurrentNote(t *testing.T) {
	resetViper(t)
	useTempWorkspace(t)

	cmd, _ := newTestCommand()
	err := showNote(cmd, nil)
	if err == nil {
		t.Fatal("showNote() expected an error")
	}
	if !strings.Contains(err.Error(), "no current note selected") {
		t.Errorf("error = %q, want current-note guidance", err)
	}
}

func TestShowNoteRejectsMissingNamedNote(t *testing.T) {
	resetViper(t)
	useTempWorkspace(t)
	if err := os.MkdirAll(defaultNotesDirName, 0755); err != nil {
		t.Fatal(err)
	}

	cmd, _ := newTestCommand()
	err := showNote(cmd, []string{"missing.md"})
	if err == nil {
		t.Fatal("showNote() expected an error")
	}
	if !strings.Contains(err.Error(), "does not exist") {
		t.Errorf("error = %q, want missing-note guidance", err)
	}
}

func TestShowNoteRejectsUnsafeFilename(t *testing.T) {
	tests := []string{"../other.md", "/tmp/other.md", "nested/other.md"}

	for _, filename := range tests {
		t.Run(filename, func(t *testing.T) {
			resetViper(t)
			useTempWorkspace(t)

			cmd, _ := newTestCommand()
			err := showNote(cmd, []string{filename})
			if err == nil {
				t.Fatal("showNote() expected an error")
			}
			if !strings.Contains(err.Error(), "single filename") {
				t.Errorf("error = %q, want filename guidance", err)
			}
		})
	}
}

func TestShowNoteRejectsDirectory(t *testing.T) {
	resetViper(t)
	useTempWorkspace(t)
	if err := os.MkdirAll(filepath.Join(defaultNotesDirName, "ideas.md"), 0755); err != nil {
		t.Fatal(err)
	}

	cmd, _ := newTestCommand()
	err := showNote(cmd, []string{"ideas.md"})
	if err == nil {
		t.Fatal("showNote() expected an error")
	}
	if !strings.Contains(err.Error(), "not a regular file") {
		t.Errorf("error = %q, want file-type guidance", err)
	}
}
