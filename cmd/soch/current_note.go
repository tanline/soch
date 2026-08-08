package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const currentNotePath = ".soch/current-note"

func setCurrentNote(notePath string) error {
	if err := validateCurrentNotePath(); err != nil {
		return err
	}

	info, err := os.Lstat(currentNotePath)
	if err == nil {
		if info.Mode()&os.ModeSymlink == 0 {
			return fmt.Errorf("%q already exists and is not a symlink; remove or rename it", currentNotePath)
		}
		if err := os.Remove(currentNotePath); err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	target, err := filepath.Rel(filepath.Dir(currentNotePath), notePath)
	if err != nil {
		return err
	}

	return os.Symlink(target, currentNotePath)
}

func validateCurrentNotePath() error {
	info, err := os.Lstat(currentNotePath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("checking current note path %q: %w", currentNotePath, err)
	}
	if info.Mode()&os.ModeSymlink == 0 {
		return fmt.Errorf("%q already exists and is not a symlink; remove or rename it", currentNotePath)
	}

	return nil
}

func getCurrentNotePath() (string, error) {
	info, err := os.Lstat(currentNotePath)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("no current note selected; run `soch new \"<title>\"` first")
	}
	if err != nil {
		return "", fmt.Errorf("checking current note path %q: %w", currentNotePath, err)
	}
	if info.Mode()&os.ModeSymlink == 0 {
		return "", fmt.Errorf("current note path %q is not a symlink; run `soch new \"<title>\"` after resolving it", currentNotePath)
	}

	target, err := os.Readlink(currentNotePath)
	if err != nil {
		return "", fmt.Errorf("reading current note path %q: %w", currentNotePath, err)
	}

	notePath := resolveLinkTarget(currentNotePath, target)
	noteInfo, err := os.Stat(notePath)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("current note target %q does not exist; run `soch new \"<title>\"` to select a note", notePath)
	}
	if err != nil {
		return "", fmt.Errorf("checking current note target %q: %w", notePath, err)
	}
	if !noteInfo.Mode().IsRegular() {
		return "", fmt.Errorf("current note target %q is not a regular file", notePath)
	}

	return notePath, nil
}
