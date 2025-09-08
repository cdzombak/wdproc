package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func MoveToTrash(filePath string, dryRun, verbose bool) error {
	if dryRun {
		fmt.Printf("[DRY RUN] Would move to trash: %s\n", filePath)
		return nil
	}

	if verbose {
		fmt.Printf("Moving to trash: %s\n", filePath)
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filePath)
	}

	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	cmd := exec.Command("osascript", "-e", fmt.Sprintf(`tell application "Finder" to delete POSIX file "%s"`, absPath))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to move file to trash: %w, output: %s", err, string(output))
	}

	if verbose {
		fmt.Printf("Successfully moved to trash: %s\n", filePath)
	}

	return nil
}