package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestSetMasterNotesDir(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)

	masterNotesDir := t.TempDir()
	cmd := &cobra.Command{}
	cmd.SetIn(strings.NewReader(masterNotesDir + "\n"))
	cmd.SetOut(&bytes.Buffer{})

	err := setMasterNotesDir(cmd)
	if err != nil {
		t.Fatalf("setMasterNotesDir() returned an error: %v", err)
	}

	if got := viper.GetString("master_notes_dir"); got != masterNotesDir {
		t.Errorf("master_notes_dir = %q, want %q", got, masterNotesDir)
	}
}

func TestSetMasterNotesDirRejectsRelativePath(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)

	cmd := &cobra.Command{}
	cmd.SetIn(strings.NewReader("notes\n"))
	cmd.SetOut(&bytes.Buffer{})

	err := setMasterNotesDir(cmd)
	if err == nil {
		t.Fatal("setMasterNotesDir() expected an error for a relative path")
	}

	if !strings.Contains(err.Error(), "must be absolute") {
		t.Errorf("error = %q, want an absolute-path error", err)
	}
}
