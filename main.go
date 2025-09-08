package main

import (
	"flag"
	"fmt"
	"log"
)

const version = "0.1.0"

func main() {
	var (
		configPath  = flag.String("config", "./config.yaml", "Path to configuration file")
		showVersion = flag.Bool("version", false, "Show version information and exit")
		verbose     = flag.Bool("verbose", false, "Enable verbose logging")
		dryRun      = flag.Bool("dry-run", true, "Preview actions without executing")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("wdproc version %s\n", version)
		return
	}

	setupLogging(*verbose)

	processor, err := NewProcessor(*configPath, *verbose, *dryRun)
	if err != nil {
		log.Fatalf("Failed to create processor: %v", err)
	}

	if err := processor.Process(); err != nil {
		log.Fatalf("Processing failed: %v", err)
	}
}

func setupLogging(verbose bool) {
	if verbose {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	} else {
		log.SetFlags(log.LstdFlags)
	}
}
