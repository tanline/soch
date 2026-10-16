package soch

import (
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestSetMasterNotesDir(t *testing.T) {
	resetViper(t)

	masterNotesDir := t.TempDir()
	cmd, _ := newTestCommand()
	cmd.SetIn(strings.NewReader(masterNotesDir + "\n"))

	err := setMasterNotesDir(cmd)
	if err != nil {
		t.Fatalf("setMasterNotesDir() returned an error: %v", err)
	}

	if got := viper.GetString("master_notes_dir"); got != masterNotesDir {
		t.Errorf("master_notes_dir = %q, want %q", got, masterNotesDir)
	}
}

func TestSetMasterNotesDirRejectsRelativePath(t *testing.T) {
	resetViper(t)

	cmd, _ := newTestCommand()
	cmd.SetIn(strings.NewReader("notes\n"))

	err := setMasterNotesDir(cmd)
	if err == nil {
		t.Fatal("setMasterNotesDir() expected an error for a relative path")
	}

	if !strings.Contains(err.Error(), "must be absolute") {
		t.Errorf("error = %q, want an absolute-path error", err)
	}
}
