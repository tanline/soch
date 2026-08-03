package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultNotesDirName       = "notes"
	projectMetadataDirName    = ".soch"
	notesDirNameKey           = "notes_dir_name"
	addNotesDirToGitignoreKey = "add_notes_dir_to_gitignore"
	masterNotesDirKey         = "master_notes_dir"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a notes directory in current workspace",
	RunE:  createNotesDir,
}

func createNotesDir(cmd *cobra.Command, args []string) error {
	notesDirName := getNotesDirName()
	if err := validateInitPrerequisites(notesDirName); err != nil {
		return err
	}

	if err := createLocalNotesDir(notesDirName); err != nil {
		return fmt.Errorf("creating notes directory: %w", err)
	}
	if err := createProjectMetadataDir(); err != nil {
		return fmt.Errorf("creating project metadata directory: %w", err)
	}

	if shouldAddNotesDirToGitignore() {
		added, err := updateGitignore(notesDirName)
		if err != nil {
			return fmt.Errorf("updating .gitignore: %w", err)
		}

		if added {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Adding %s directory to .gitignore\n", notesDirName); err != nil {
				return err
			}
		}
	}

	added, err := updateGitignore(projectMetadataDirName)
	if err != nil {
		return fmt.Errorf("updating .gitignore: %w", err)
	}
	if added {
		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Adding %s directory to .gitignore\n", projectMetadataDirName); err != nil {
			return err
		}
	}

	linked, err := linkNotesDir(notesDirName)
	if err != nil {
		return fmt.Errorf("linking notes directory: %w", err)
	}

	if linked {
		if _, err := fmt.Fprintln(cmd.OutOrStdout(), "Linked notes directory to master notes"); err != nil {
			return err
		}
	}

	_, err = fmt.Fprintln(cmd.OutOrStdout(), "Notes directory is ready")
	return err
}

func validateInitPrerequisites(notesDirName string) error {
	if err := validateNotesDirName(notesDirName); err != nil {
		return err
	}

	workspaceDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot determine the current workspace: %w", err)
	}

	masterNotesDir := strings.TrimSpace(viper.GetString(masterNotesDirKey))
	if masterNotesDir == "" {
		return fmt.Errorf("master notes directory is not configured; run `soch setup` first")
	}

	info, err := os.Stat(masterNotesDir)
	if err != nil {
		return fmt.Errorf("cannot access master notes directory %q: %w; check the path in your soch config", masterNotesDir, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("master notes path %q is not a directory; update your soch config", masterNotesDir)
	}

	linkPath := filepath.Join(masterNotesDir, filepath.Base(workspaceDir))
	expectedTarget := filepath.Join(workspaceDir, notesDirName)
	if _, err := validateExistingLink(linkPath, expectedTarget); err != nil {
		return err
	}

	metadataInfo, err := os.Stat(projectMetadataDirName)
	if err == nil {
		if !metadataInfo.IsDir() {
			return fmt.Errorf("project metadata path %q already exists and is not a directory; remove or rename it before running `soch init`", projectMetadataDirName)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("cannot inspect project metadata path %q: %w", projectMetadataDirName, err)
	}

	return nil
}

func createLocalNotesDir(notesDirName string) error {
	return os.MkdirAll(notesDirName, 0755)
}

func createProjectMetadataDir() error {
	return os.MkdirAll(projectMetadataDirName, 0755)
}

func linkNotesDir(notesDirName string) (bool, error) {
	masterNotesDir := strings.TrimSpace(viper.GetString(masterNotesDirKey))
	if masterNotesDir == "" {
		return false, fmt.Errorf("master notes directory is not configured")
	}

	info, err := os.Stat(masterNotesDir)
	if err != nil {
		return false, fmt.Errorf("accessing master notes directory %q: %w", masterNotesDir, err)
	}
	if !info.IsDir() {
		return false, fmt.Errorf("master notes path %q is not a directory", masterNotesDir)
	}

	workspaceDir, err := os.Getwd()
	if err != nil {
		return false, fmt.Errorf("getting current workspace: %w", err)
	}

	linkPath := filepath.Join(masterNotesDir, filepath.Base(workspaceDir))
	notesPath := filepath.Join(workspaceDir, notesDirName)

	exists, err := validateExistingLink(linkPath, notesPath)
	if err != nil {
		return false, err
	}
	if exists {
		return false, nil
	}

	if err := os.Symlink(notesPath, linkPath); err != nil {
		return false, err
	}

	return true, nil
}

func validateExistingLink(linkPath, expectedTarget string) (bool, error) {
	linkInfo, err := os.Lstat(linkPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("cannot inspect link path %q: %w", linkPath, err)
	}
	if linkInfo.Mode()&os.ModeSymlink == 0 {
		return false, fmt.Errorf("link path %q already exists; remove or rename it before running `soch init`", linkPath)
	}

	existingTarget, err := os.Readlink(linkPath)
	if err != nil {
		return false, fmt.Errorf("cannot inspect existing link %q: %w", linkPath, err)
	}

	resolvedExistingTarget := resolveLinkTarget(linkPath, existingTarget)
	resolvedExpectedTarget := filepath.Clean(expectedTarget)

	if _, err := os.Stat(expectedTarget); err == nil {
		resolvedExistingTarget, err = filepath.EvalSymlinks(resolvedExistingTarget)
		if err != nil {
			return false, fmt.Errorf("resolving existing symlink target: %w", err)
		}

		resolvedExpectedTarget, err = filepath.EvalSymlinks(expectedTarget)
		if err != nil {
			return false, fmt.Errorf("resolving expected symlink target: %w", err)
		}

		if resolvedExistingTarget == resolvedExpectedTarget {
			return true, nil
		}
	} else if !os.IsNotExist(err) {
		return false, fmt.Errorf("checking expected symlink target: %w", err)
	}

	if resolvedExistingTarget != resolvedExpectedTarget {
		return false, fmt.Errorf("link path %q points to %q instead of %q; remove or rename it before running `soch init`", linkPath, existingTarget, expectedTarget)
	}

	return true, nil
}

func resolveLinkTarget(linkPath, target string) string {
	if filepath.IsAbs(target) {
		return filepath.Clean(target)
	}

	return filepath.Clean(filepath.Join(filepath.Dir(linkPath), target))
}

func shouldAddNotesDirToGitignore() bool {
	if !viper.IsSet(addNotesDirToGitignoreKey) {
		return true
	}

	return viper.GetBool(addNotesDirToGitignoreKey)
}

func updateGitignore(entryName string) (bool, error) {
	const gitignoreFileName = ".gitignore"

	contents, err := os.ReadFile(gitignoreFileName)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	gitignoreEntry := entryName + "/"
	targetEntry := strings.Trim(entryName, "/")
	for line := range strings.SplitSeq(strings.ReplaceAll(string(contents), "\r\n", "\n"), "\n") {
		entry := strings.Trim(strings.TrimSpace(line), "/")
		if entry == targetEntry {
			return false, nil
		}
	}

	updatedContents := string(contents)
	if len(updatedContents) > 0 && !strings.HasSuffix(updatedContents, "\n") {
		updatedContents += "\n"
	}
	updatedContents += gitignoreEntry + "\n"

	if err := os.WriteFile(gitignoreFileName, []byte(updatedContents), 0644); err != nil {
		return false, err
	}

	return true, nil
}

func getNotesDirName() string {
	name := strings.TrimSpace(viper.GetString(notesDirNameKey))

	if name == "" {
		return defaultNotesDirName
	}

	return name
}

func validateNotesDirName(notesDirName string) error {
	if notesDirName == "." || notesDirName == ".." ||
		filepath.IsAbs(notesDirName) || filepath.Base(notesDirName) != notesDirName {
		return fmt.Errorf("invalid notes directory name %q; use a single relative directory name such as `notes`", notesDirName)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)
}
