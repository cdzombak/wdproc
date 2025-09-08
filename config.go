package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Version      int        `yaml:"version"`
	Paths        Paths      `yaml:"paths"`
	Processing   Processing `yaml:"processing"`
	Integrations Integrations `yaml:"integrations"`
}

type Paths struct {
	Journal string `yaml:"journal"`
	Tasks   string `yaml:"tasks"`
}

type Processing struct {
	JournalMinAge string `yaml:"journal_min_age"`
	TasksMinAge   string `yaml:"tasks_min_age"`
}

type Integrations struct {
	DayOne DayOneConfig `yaml:"day_one"`
	Things ThingsConfig `yaml:"things"`
}

type DayOneConfig struct {
	Enabled bool     `yaml:"enabled"`
	Journal string   `yaml:"journal"`
	Tags    []string `yaml:"tags"`
}

type ThingsConfig struct {
	Enabled bool     `yaml:"enabled"`
	Tags    []string `yaml:"tags"`
}


func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	if err := config.ExpandPaths(); err != nil {
		return nil, fmt.Errorf("failed to expand paths: %w", err)
	}

	return &config, nil
}

func (c *Config) Validate() error {
	if c.Version != 1 {
		return fmt.Errorf("unsupported config version: %d", c.Version)
	}

	if c.Paths.Journal == "" {
		return fmt.Errorf("paths.journal is required")
	}
	if c.Paths.Tasks == "" {
		return fmt.Errorf("paths.tasks is required")
	}

	if _, err := time.ParseDuration(c.Processing.JournalMinAge); err != nil {
		return fmt.Errorf("invalid processing.journal_min_age: %w", err)
	}
	if _, err := time.ParseDuration(c.Processing.TasksMinAge); err != nil {
		return fmt.Errorf("invalid processing.tasks_min_age: %w", err)
	}

	if c.Integrations.DayOne.Enabled && c.Integrations.DayOne.Journal == "" {
		return fmt.Errorf("integrations.day_one.journal is required when Day One is enabled")
	}

	return nil
}

func (c *Config) ExpandPaths() error {
	journalPath, err := expandPath(c.Paths.Journal)
	if err != nil {
		return fmt.Errorf("failed to expand journal path: %w", err)
	}
	c.Paths.Journal = journalPath

	tasksPath, err := expandPath(c.Paths.Tasks)
	if err != nil {
		return fmt.Errorf("failed to expand tasks path: %w", err)
	}
	c.Paths.Tasks = tasksPath

	return nil
}

func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(homeDir, path[2:]), nil
	}
	return filepath.Abs(path)
}

func (c *Config) GetJournalMinAge() (time.Duration, error) {
	return time.ParseDuration(c.Processing.JournalMinAge)
}

func (c *Config) GetTasksMinAge() (time.Duration, error) {
	return time.ParseDuration(c.Processing.TasksMinAge)
}