package soch

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestDefaultConfigFilePath(t *testing.T) {
	homeDir := filepath.Join("testdata", "home")
	want := filepath.Join(homeDir, ".config", "soch", "config.yaml")

	if got := defaultConfigFilePath(homeDir); got != want {
		t.Errorf("defaultConfigFilePath() = %q, want %q", got, want)
	}
}

func TestGetConfigFilePathUsesFlag(t *testing.T) {
	cmd := &cobra.Command{}
	configPath := filepath.Join(t.TempDir(), "config.yaml")
	cmd.Flags().String(configFlagName, configPath, "")

	got, err := getConfigFilePath(cmd)
	if err != nil {
		t.Fatalf("getConfigFilePath() returned an error: %v", err)
	}
	if got != configPath {
		t.Errorf("getConfigFilePath() = %q, want %q", got, configPath)
	}
}

func TestGetConfigFilePathUsesDefault(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(configFlagName, "", "")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("getting user home directory: %v", err)
	}

	got, err := getConfigFilePath(cmd)
	if err != nil {
		t.Fatalf("getConfigFilePath() returned an error: %v", err)
	}
	want := defaultConfigFilePath(homeDir)
	if got != want {
		t.Errorf("getConfigFilePath() = %q, want %q", got, want)
	}
}

func TestLoadConfigFile(t *testing.T) {
	resetViper(t)

	configPath := filepath.Join("testdata", "config.yaml")

	if err := loadConfigFile(configPath); err != nil {
		t.Fatalf("loadConfigFile() returned an error: %v", err)
	}

	if got := viper.GetString(notesDirNameKey); got != "project-notes" {
		t.Errorf("notes_dir_name = %q, want %q", got, "project-notes")
	}
}

func TestLoadConfigFileAllowsMissingFile(t *testing.T) {
	resetViper(t)

	configPath := filepath.Join(t.TempDir(), "missing.yaml")
	if err := loadConfigFile(configPath); err != nil {
		t.Fatalf("loadConfigFile() returned an error: %v", err)
	}
}
