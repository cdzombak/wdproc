package main

import (
	"fmt"
	"log"
	"time"
)

type Processor struct {
	config        *Config
	dayOne        *DayOneIntegration
	things        *ThingsIntegration
	verbose       bool
	dryRun        bool
}

func NewProcessor(configPath string, verbose, dryRun bool) (*Processor, error) {
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}


	dayOne := NewDayOneIntegration(config.Integrations.DayOne, dryRun, verbose)
	things := NewThingsIntegration(config.Integrations.Things, dryRun, verbose)

	return &Processor{
		config:  config,
		dayOne:  dayOne,
		things:  things,
		verbose: verbose,
		dryRun:  dryRun,
	}, nil
}

func (p *Processor) Process() error {
	log.Println("Starting processing run")
	now := time.Now()

	if p.config.Integrations.DayOne.Enabled {
		if err := p.processJournalFiles(now); err != nil {
			return fmt.Errorf("failed to process journal files: %w", err)
		}
	}

	if p.config.Integrations.Things.Enabled {
		if err := p.processTaskFiles(now); err != nil {
			return fmt.Errorf("failed to process task files: %w", err)
		}
	}

	log.Println("Processing run completed")
	return nil
}

func (p *Processor) processJournalFiles(now time.Time) error {
	minAge, err := p.config.GetJournalMinAge()
	if err != nil {
		return fmt.Errorf("invalid journal min age: %w", err)
	}

	files, err := ScanDirectory(p.config.Paths.Journal, minAge, now)
	if err != nil {
		return fmt.Errorf("failed to scan journal directory: %w", err)
	}

	log.Printf("Found %d eligible journal files", len(files))

	for _, file := range files {
		if err := p.processJournalFile(file); err != nil {
			log.Printf("WARNING: Failed to process journal file %s: %v", file.Name, err)
			continue
		}
	}

	return nil
}

func (p *Processor) processTaskFiles(now time.Time) error {
	minAge, err := p.config.GetTasksMinAge()
	if err != nil {
		return fmt.Errorf("invalid tasks min age: %w", err)
	}

	files, err := ScanDirectory(p.config.Paths.Tasks, minAge, now)
	if err != nil {
		return fmt.Errorf("failed to scan tasks directory: %w", err)
	}

	log.Printf("Found %d eligible task files", len(files))

	for _, file := range files {
		if err := p.processTaskFile(file); err != nil {
			log.Printf("WARNING: Failed to process task file %s: %v", file.Name, err)
			continue
		}
	}

	return nil
}

func (p *Processor) processJournalFile(file FileInfo) error {
	if p.verbose {
		log.Printf("Processing journal file: %s (age: %v)", file.Name, file.Age.Truncate(time.Minute))
	}

	if err := p.dayOne.CreateEntry(file); err != nil {
		return fmt.Errorf("failed to create Day One entry: %w", err)
	}

	if err := MoveToTrash(file.Path, p.dryRun, p.verbose); err != nil {
		return fmt.Errorf("failed to move file to trash: %w", err)
	}

	return nil
}

func (p *Processor) processTaskFile(file FileInfo) error {
	if p.verbose {
		log.Printf("Processing task file: %s (age: %v)", file.Name, file.Age.Truncate(time.Minute))
	}

	if err := p.things.CreateTask(file); err != nil {
		return fmt.Errorf("failed to create Things task: %w", err)
	}

	if err := MoveToTrash(file.Path, p.dryRun, p.verbose); err != nil {
		return fmt.Errorf("failed to move file to trash: %w", err)
	}

	return nil
}