package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const thoughtFlagName = "thought"

var addCmd = &cobra.Command{
	Use:   "add [thought]",
	Short: "Add a thought to the current note",
	Args:  cobra.MaximumNArgs(1),
	RunE:  addThought,
}

func addThought(cmd *cobra.Command, args []string) error {
	thought, err := cmd.Flags().GetString(thoughtFlagName)
	if err != nil {
		return err
	}
	if len(args) > 0 && cmd.Flags().Changed(thoughtFlagName) {
		return fmt.Errorf("provide a thought as an argument or with --%s, not both", thoughtFlagName)
	}

	return addThoughtAt(cmd, args, thought, time.Now())
}

func addThoughtAt(cmd *cobra.Command, args []string, thoughtFlag string, now time.Time) error {
	notePath, err := getCurrentNotePath()
	if err != nil {
		return err
	}

	thoughts, err := readThoughts(cmd, args, thoughtFlag)
	if err != nil {
		return err
	}

	if err := appendThought(notePath, thoughts, now); err != nil {
		return fmt.Errorf("adding thought to %q: %w", notePath, err)
	}

	_, err = fmt.Fprintf(cmd.OutOrStdout(), "Added thought to %s\n", notePath)
	return err
}

func readThoughts(cmd *cobra.Command, args []string, thoughtFlag string) (string, error) {
	if len(args) > 0 && strings.TrimSpace(thoughtFlag) != "" {
		return "", fmt.Errorf("provide a thought as an argument or with --%s, not both", thoughtFlagName)
	}

	if len(args) > 0 {
		return validateThought(args[0])
	}

	if strings.TrimSpace(thoughtFlag) != "" {
		return validateThought(thoughtFlag)
	}

	if isInteractiveInput(cmd.InOrStdin()) {
		if _, err := fmt.Fprint(cmd.ErrOrStderr(), "Thought: "); err != nil {
			return "", err
		}
	}

	input, err := io.ReadAll(cmd.InOrStdin())
	if err != nil {
		return "", fmt.Errorf("reading thought: %w", err)
	}

	return validateThought(string(input))
}

func validateThought(thought string) (string, error) {
	thought = strings.TrimSpace(thought)
	if thought == "" {
		return "", fmt.Errorf("thought cannot be empty")
	}

	return thought, nil
}

func isInteractiveInput(reader io.Reader) bool {
	file, ok := reader.(*os.File)
	if !ok {
		return false
	}

	info, err := file.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
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

func appendThought(notePath, newThought string, now time.Time) error {
	var entry strings.Builder
	lines := strings.Split(strings.ReplaceAll(newThought, "\r\n", "\n"), "\n")

	fmt.Fprintf(&entry, "- %s — %s\n", now.Format(time.RFC3339), strings.TrimSpace(lines[0]))
	for _, line := range lines[1:] {
		fmt.Fprintf(&entry, "  %s\n", strings.TrimSpace(line))
	}

	file, err := os.OpenFile(notePath, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(entry.String())
	return err
}

func init() {
	addCmd.Flags().String(thoughtFlagName, "", "Thought to add to the current note")
	rootCmd.AddCommand(addCmd)
}
