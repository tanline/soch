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

func TestNewRootCommandRegistersCommandsAndFlags(t *testing.T) {
	first := newRootCommand()
	second := newRootCommand()
	if first == second {
		t.Fatal("newRootCommand() returned the same command twice")
	}

	for _, root := range []*cobra.Command{first, second} {
		if root.PersistentFlags().Lookup(configFlagName) == nil {
			t.Error("root command is missing --config")
		}

		for _, name := range []string{"setup", "init", "new", "add", "show"} {
			command, _, err := root.Find([]string{name})
			if err != nil || command == nil || command == root || command.Name() != name {
				t.Errorf("root command is missing %q: command = %v, error = %v", name, command, err)
			}
			if name == "add" && command != nil && command.Flags().Lookup(thoughtFlagName) == nil {
				t.Error("add command is missing --thought")
			}
		}
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
