# WriterDeck Processor (wdproc) Specification

## Overview

`wdproc` is a Go program that monitors specific WriterDeck folders for files and automatically processes them by creating entries in Day One journal and tasks in Things 3, then moving the processed files to the macOS Trash.

## Core Functionality

### File Processing

1. **Journal Folder Processing** (`~/WriterDeck/Journal`)
   - Monitors for files older than 4 hours (configurable)
   - Creates a Day One journal entry for each file:
     - Entry date: File's creation timestamp
     - Entry title: File name without extension
     - Entry content: File contents
   - Moves processed files to macOS Trash

2. **Tasks Folder Processing** (`~/WriterDeck/Tasks`)
   - Monitors for files older than 1 hour (configurable)
   - Creates a Things 3 task for each file:
     - Task title: File name without extension
     - Task notes: File contents
     - Task location: Inbox
   - Moves processed files to macOS Trash

### File Age Calculation
- File age is determined by the file's creation time
- Only files exceeding the configured minimum age are processed

## Configuration

### Configuration File Format

The program accepts a YAML configuration file path via command-line flag:
```bash
wdproc --config=/path/to/config.yaml
```

### Configuration Schema

```yaml
# config.yaml
version: 1

paths:
  journal: "~/WriterDeck/Journal"    # Path to journal files
  tasks: "~/WriterDeck/Tasks"        # Path to task files

processing:
  journal_min_age: "4h"               # Minimum age for journal files (Go duration format)
  tasks_min_age: "1h"                 # Minimum age for task files (Go duration format)

integrations:
  day_one:
    enabled: true                     # Enable/disable Day One integration
    journal: "Default"                # Target journal name (required if enabled)
    tags: []                          # Optional: Tags to add to entries

  things:
    enabled: true                     # Enable/disable Things integration
    tags: []                          # Optional: Tags to add to tasks
```

## Integration Details

### Day One Integration

The program will create journal entries using the Day One CLI tool located at `/usr/local/bin/dayone2`.

#### CLI Command Structure

```bash
/usr/local/bin/dayone2 new [entry content] \
  --date='[formatted date]' \
  --journal='[journal name]' \
  --tags [tag1] [tag2] \
  [additional options]
```

#### Entry Creation Details

1. **Entry Title**: The filename (without extension) will be included as the first line of the entry content
2. **Entry Date**: Use the file's creation timestamp, formatted as `yyyy-mm-dd hh:mm:ss`
3. **Entry Content**: The complete file contents following the title line
4. **Journal Selection**: Use the journal name from `integrations.day_one.journal` configuration
5. **Tags**: Apply any configured tags to the entry

#### Example Command

For a file named "Morning Thoughts.md" created on 2024-01-15 at 09:30:00, with configuration `journal: "Personal"` and `tags: ["Writing", "Journal"]`:
```bash
/usr/local/bin/dayone2 new "Morning Thoughts

[file contents here]" \
  --date='2024-01-15 09:30:00' \
  --journal='Personal' \
  --tags Writing Journal
```

### Things 3 Integration

The program will create tasks using the Things URL scheme:

```
things:///add?title=[URL-encoded title]&notes=[URL-encoded content]
```

Tasks will be created in the Inbox by default.

## Processing Logic

### Main Processing Loop

1. Load and validate configuration
2. For each enabled integration:
   - Scan the appropriate folder
   - Filter files by minimum age
   - Process each eligible file:
     - Read file contents
     - Create integration entry (Day One or Things)
     - Move file to Trash if successful
   - Log processing results

### Error Handling

- **File Read Errors**: Log error, skip file, continue processing
- **Integration Errors**: Log error, do not move file to trash
- **Trash Move Errors**: Log error, leave file in place
- **Configuration Errors**: Exit with error message

### File Filtering

Files are processed if:
- They are regular files (not directories)
- They are not hidden files (don't start with `.`)
- They exceed the minimum age threshold
- They have read permissions

## Command-Line Interface

```bash
Usage: wdproc [OPTIONS]

Options:
  --config PATH     Path to configuration file (required)
  --version         Show version information
  --help            Show help message
  --verbose         Enable verbose logging
  --dry-run         Preview actions without executing
```

## Logging

The program should log:
- Start/stop of processing runs
- Number of files found in each folder
- Each file processed (name, age, destination)
- Integration responses/errors
- Files moved to trash
- Configuration validation errors

Log levels:
- `INFO`: Normal operations
- `WARN`: Recoverable errors (e.g., single file processing failure)
- `ERROR`: Fatal errors requiring attention

## Platform Requirements

- **OS**: macOS (required for Trash integration, URL schemes, and Day One CLI)
- **Go Version**: 1.21 or higher
- **External Tools**:
  - Day One CLI: `/usr/local/bin/dayone2` (required for Day One integration)
  - Day One.app installed
- **Dependencies**:
  - YAML parsing library (e.g., `gopkg.in/yaml.v3`)
  - macOS Trash integration
  - URL encoding utilities for Things integration
  - Command execution utilities for CLI interaction

## Security Considerations

- Configuration file should have appropriate permissions (readable by user only)
- File contents should be properly escaped when passing to shell commands (Day One CLI)
- File contents should be properly URL-encoded before passing to URL schemes (Things)
- The program should not follow symbolic links to prevent directory traversal
- Shell command execution should use proper argument escaping to prevent injection

## Future Enhancements

1. **Watch Mode**: Continuously monitor folders instead of single-run processing
2. **Custom Templates**: Allow templating for entry/task creation
3. **Multiple Folder Support**: Process arbitrary number of folder pairs
4. **Backup Before Delete**: Option to backup files before moving to trash
5. **Things Projects**: Support for adding tasks to specific projects/areas
6. **File Type Filtering**: Process only specific file extensions
7. **Scheduling**: Built-in scheduling without external cron dependency
8. **Day One Attachments**: Support for adding file attachments to entries
9. **Day One Location**: Add location data to journal entries

## Testing Strategy

1. **Unit Tests**:
   - Configuration parsing and validation
   - File age calculation
   - URL encoding functions
   - File filtering logic

2. **Integration Tests**:
   - Mock file system operations
   - Mock Day One CLI command execution
   - Mock URL scheme calls (Things)
   - Trash operation simulation

3. **Manual Testing**:
   - Test with actual Day One and Things applications
   - Verify trash operations
   - Test with various file types and encodings
