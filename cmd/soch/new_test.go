package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestCreateNewNoteCreatesCollectionAndSelectsIt(t *testing.T) {
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

	if err := createNewNoteAt(cmd, []string{"Ideas & Plans"}, now); err != nil {
		t.Fatalf("createNewNoteAt() returned an error: %v", err)
	}

	notePath := filepath.Join(defaultNotesDirName, "2026-08-02-ideas-plans.md")
	contents, err := os.ReadFile(notePath)
	if err != nil {
		t.Fatalf("reading collection: %v", err)
	}
	if string(contents) != "# Ideas & Plans\n\n" {
		t.Errorf("collection = %q, want %q", contents, "# Ideas & Plans\n\n")
	}

	target, err := os.Readlink(currentNotePath)
	if err != nil {
		t.Fatalf("reading current note link: %v", err)
	}
	if target != filepath.Join("..", notePath) {
		t.Errorf("current note target = %q, want %q", target, filepath.Join("..", notePath))
	}

	if !strings.Contains(output.String(), "Created collection: "+notePath) {
		t.Error("output did not report the created collection")
	}
}

func TestCreateNewNoteReusesExistingCollection(t *testing.T) {
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
	if err := createNewNoteAt(cmd, []string{"Ideas"}, now); err != nil {
		t.Fatalf("createNewNoteAt() returned an error: %v", err)
	}

	contents, err := os.ReadFile(notePath)
	if err != nil {
		t.Fatal(err)
	}
	if string(contents) != existingContents {
		t.Errorf("existing collection was modified: %q", contents)
	}
	if !strings.Contains(output.String(), "Using existing collection: "+notePath) {
		t.Error("output did not report reusing the collection")
	}
}

func TestCreateNewNoteRequiresInitializedProject(t *testing.T) {
	resetViper(t)
	useTempWorkspace(t)

	cmd, _ := newTestCommand()
	err := createNewNoteAt(cmd, []string{"Ideas"}, time.Now())
	if err == nil {
		t.Fatal("createNewNoteAt() expected an initialization error")
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

func TestCreateNewNoteRejectsConflictingCurrentPointer(t *testing.T) {
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
	err := createNewNoteAt(cmd, []string{"Ideas"}, time.Now())
	if err == nil {
		t.Fatal("createNewNoteAt() expected a current pointer conflict")
	}
	if !strings.Contains(err.Error(), "not a symlink") {
		t.Errorf("error = %q, want symlink guidance", err)
	}
}
