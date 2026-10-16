package soch

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestCreateNewNoteAndSetCurrentCreatesNote(t *testing.T) {
	resetViper(t)
	useTempWorkspace(t)
	if err := os.MkdirAll(defaultNotesDirName, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(projectMetadataDirName, 0755); err != nil {
		t.Fatal(err)
	}

	now := time.Date(2026, time.August, 2, 14, 32, 5, 0, time.FixedZone("EDT", -4*60*60))
	cmd, output := newTestCommand()

	if err := createNewNoteAndSetCurrentAt(cmd, []string{"Ideas & Plans"}, now); err != nil {
		t.Fatalf("createNewNoteAndSetCurrentAt() returned an error: %v", err)
	}

	notePath := filepath.Join(defaultNotesDirName, "2026-08-02-ideas-plans.md")
	contents, err := os.ReadFile(notePath)
	if err != nil {
		t.Fatalf("reading note: %v", err)
	}
	if string(contents) != "# Ideas & Plans\n\n" {
		t.Errorf("note = %q, want %q", contents, "# Ideas & Plans\n\n")
	}

	target, err := os.Readlink(currentNotePath)
	if err != nil {
		t.Fatalf("reading current note link: %v", err)
	}
	if target != filepath.Join("..", notePath) {
		t.Errorf("current note target = %q, want %q", target, filepath.Join("..", notePath))
	}

	if !strings.Contains(output.String(), "Created note: "+notePath) {
		t.Error("output did not report the created note")
	}
}

func TestCreateNewNoteAndSetCurrentReusesExistingNote(t *testing.T) {
	resetViper(t)
	useTempWorkspace(t)
	if err := os.MkdirAll(defaultNotesDirName, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(projectMetadataDirName, 0755); err != nil {
		t.Fatal(err)
	}

	notePath := filepath.Join(defaultNotesDirName, "2026-08-02-ideas.md")
	existingContents := "# Ideas\n\n- Existing thought\n"
	if err := os.WriteFile(notePath, []byte(existingContents), 0644); err != nil {
		t.Fatal(err)
	}

	cmd, output := newTestCommand()
	now := time.Date(2026, time.August, 2, 14, 32, 5, 0, time.UTC)
	if err := createNewNoteAndSetCurrentAt(cmd, []string{"Ideas"}, now); err != nil {
		t.Fatalf("createNewNoteAndSetCurrentAt() returned an error: %v", err)
	}

	contents, err := os.ReadFile(notePath)
	if err != nil {
		t.Fatal(err)
	}
	if string(contents) != existingContents {
		t.Errorf("existing note was modified: %q", contents)
	}
	if !strings.Contains(output.String(), "Using existing note: "+notePath) {
		t.Error("output did not report reusing the note")
	}
}

func TestCreateNewNoteAndSetCurrentRequiresInitializedProject(t *testing.T) {
	resetViper(t)
	useTempWorkspace(t)

	cmd, _ := newTestCommand()
	err := createNewNoteAndSetCurrentAt(cmd, []string{"Ideas"}, time.Now())
	if err == nil {
		t.Fatal("createNewNoteAndSetCurrentAt() expected an initialization error")
	}
	if !strings.Contains(err.Error(), "run `soch init`") {
		t.Errorf("error = %q, want initialization guidance", err)
	}
}

func TestSlugifyTitle(t *testing.T) {
	tests := []struct {
		title string
		want  string
	}{
		{title: "Ideas & Plans", want: "ideas-plans"},
		{title: "  Daily   Thoughts  ", want: "daily-thoughts"},
		{title: "!!!", want: ""},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			if got := slugifyTitle(test.title); got != test.want {
				t.Errorf("slugifyTitle() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestCreateNewNoteAndSetCurrentRejectsConflictingCurrentPointer(t *testing.T) {
	resetViper(t)
	useTempWorkspace(t)
	if err := os.MkdirAll(defaultNotesDirName, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(projectMetadataDirName, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(currentNotePath, []byte("not a symlink"), 0644); err != nil {
		t.Fatal(err)
	}

	cmd, _ := newTestCommand()
	err := createNewNoteAndSetCurrentAt(cmd, []string{"Ideas"}, time.Now())
	if err == nil {
		t.Fatal("createNewNoteAndSetCurrentAt() expected a current pointer conflict")
	}
	if !strings.Contains(err.Error(), "not a symlink") {
		t.Errorf("error = %q, want symlink guidance", err)
	}
}

func TestCreateNewNoteDoesNotSetCurrentNote(t *testing.T) {
	resetViper(t)
	useTempWorkspace(t)
	if err := os.MkdirAll(defaultNotesDirName, 0755); err != nil {
		t.Fatal(err)
	}

	now := time.Date(2026, time.August, 2, 14, 32, 5, 0, time.UTC)
	notePath, created, err := createNewNote(defaultNotesDirName, "Ideas", now)
	if err != nil {
		t.Fatalf("createNewNote() returned an error: %v", err)
	}
	if !created {
		t.Error("createNewNote() created = false, want true")
	}

	wantPath := filepath.Join(defaultNotesDirName, "2026-08-02-ideas.md")
	if notePath != wantPath {
		t.Errorf("note path = %q, want %q", notePath, wantPath)
	}
	if _, err := os.Stat(notePath); err != nil {
		t.Errorf("created note does not exist: %v", err)
	}
	if _, err := os.Lstat(currentNotePath); !os.IsNotExist(err) {
		t.Errorf("current note path exists after createNewNote(): %v", err)
	}
}
