package main

import (
	"fmt"
	"net/url"
	"os/exec"
	"strings"
)

type ThingsIntegration struct {
	config  ThingsConfig
	dryRun  bool
	verbose bool
}

func NewThingsIntegration(config ThingsConfig, dryRun, verbose bool) *ThingsIntegration {
	return &ThingsIntegration{
		config:  config,
		dryRun:  dryRun,
		verbose: verbose,
	}
}

func (t *ThingsIntegration) CreateTask(file FileInfo) error {
	if !t.config.Enabled {
		return fmt.Errorf("Things integration is disabled")
	}

	thingsURL, err := t.buildThingsURL(file)
	if err != nil {
		return fmt.Errorf("failed to build Things URL: %w", err)
	}

	if t.verbose {
		fmt.Printf("Things URL: %s\n", thingsURL)
	}

	if t.dryRun {
		fmt.Printf("[DRY RUN] Would create Things task for: %s\n", file.Name)
		return nil
	}

	cmd := exec.Command("open", thingsURL)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to open Things URL: %w, output: %s", err, string(output))
	}

	if t.verbose {
		fmt.Printf("Things task created for: %s\n", file.Name)
	}

	return nil
}

func (t *ThingsIntegration) buildThingsURL(file FileInfo) (string, error) {
	baseURL := "things:///add"
	params := url.Values{}
	
	params.Set("title", file.NameNoExt)
	
	if strings.TrimSpace(file.Content) != "" {
		params.Set("notes", file.Content)
	}
	
	if len(t.config.Tags) > 0 {
		tagStr := strings.Join(t.config.Tags, ",")
		params.Set("tags", tagStr)
	}
	
	fullURL := baseURL + "?" + params.Encode()
	return fullURL, nil
}