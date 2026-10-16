package soch

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configFileName = "config"
	configFileType = "yaml"
	configFlagName = "config"
)

func defaultConfigFilePath(homeDir string) string {
	return filepath.Join(homeDir, ".config", "soch", configFileName+"."+configFileType)
}

func getConfigFilePath(cmd *cobra.Command) (string, error) {
	configPath, err := cmd.Flags().GetString(configFlagName)
	if err != nil {
		return "", fmt.Errorf("reading --%s flag: %w", configFlagName, err)
	}
	if configPath = strings.TrimSpace(configPath); configPath != "" {
		return configPath, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("finding user home directory: %w", err)
	}

	return defaultConfigFilePath(homeDir), nil
}

func loadConfig(cmd *cobra.Command) error {
	configPath, err := getConfigFilePath(cmd)
	if err != nil {
		return err
	}

	return loadConfigFile(configPath)
}

func loadConfigFile(configPath string) error {
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		if _, ok := errors.AsType[viper.ConfigFileNotFoundError](err); ok {
			return nil
		}

		return fmt.Errorf("reading config: %w", err)
	}

	return nil
}

func initSochConfig(cmd *cobra.Command) error {
	configPath, err := getConfigFilePath(cmd)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return err
	}

	viper.SetConfigFile(configPath)
	_, err = fmt.Fprintln(cmd.OutOrStdout(), "The soch config file is", configPath)
	return err
}
