package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func resetViper(t *testing.T) {
	t.Helper()

	viper.Reset()
	t.Cleanup(viper.Reset)
}

func useTempWorkspace(t *testing.T) string {
	t.Helper()

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

	return workspaceDir
}

func newTestCommand() (*cobra.Command, *bytes.Buffer) {
	var output bytes.Buffer

	cmd := &cobra.Command{}
	cmd.SetOut(&output)

	return cmd, &output
}
