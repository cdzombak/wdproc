package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const dayOneCLIPath = "/usr/local/bin/dayone2"

type DayOneIntegration struct {
	config  DayOneConfig
	dryRun  bool
	verbose bool
}

func NewDayOneIntegration(config DayOneConfig, dryRun, verbose bool) *DayOneIntegration {
	return &DayOneIntegration{
		config:  config,
		dryRun:  dryRun,
		verbose: verbose,
	}
}

func (d *DayOneIntegration) CreateEntry(file FileInfo) error {
	if !d.config.Enabled {
		//goland:noinspection GoErrorStringFormat
		return fmt.Errorf("Day One integration is disabled")
	}

	entryContent := d.buildEntryContent(file)
	dateStr := file.CreatedTime.Format("2006-01-02 15:04:05")

	args := []string{
		"new",
		entryContent,
		"--date=" + dateStr,
		"--journal=" + d.config.Journal,
	}

	if len(d.config.Tags) > 0 {
		args = append(args, "--tags")
		args = append(args, d.config.Tags...)
	}

	if d.verbose {
		fmt.Printf("Day One command: %s %s\n", dayOneCLIPath, strings.Join(args, " "))
	}

	if d.dryRun {
		fmt.Printf("[DRY RUN] Would create Day One entry for: %s\n", file.Name)
		return nil
	}

	if err := d.checkCLIExists(); err != nil {
		return err
	}

	cmd := exec.Command(dayOneCLIPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		//goland:noinspection GoErrorStringFormat
		return fmt.Errorf("Day One CLI failed: %w, output: %s", err, string(output))
	}

	if d.verbose {
		fmt.Printf("Day One entry created for: %s\n", file.Name)
	}

	return nil
}

func (d *DayOneIntegration) buildEntryContent(file FileInfo) string {
	var content strings.Builder

	content.WriteString(file.NameNoExt)
	content.WriteString("\n\n")
	content.WriteString(file.Content)

	return content.String()
}

func (d *DayOneIntegration) checkCLIExists() error {
	if _, err := os.Stat(dayOneCLIPath); os.IsNotExist(err) {
		//goland:noinspection GoErrorStringFormat
		return fmt.Errorf("Day One CLI not found at %s", dayOneCLIPath)
	}
	return nil
}
