# WriterDeck Processor (wdproc)

Monitors folders synced from a [WriterDeck](https://writerdeckos.com) and automatically converts new files into entries in Day One or tasks in Things.

I use [Syncthing](https://syncthing.net) to sync files from WriterDeck to my Mac.

## Installation

### macOS via Homebrew

```shell
brew install cdzombak/oss/wdproc
```

### Manual installation from build artifacts

Pre-built binaries for macOS on various architectures are downloadable from each [GitHub Release](https://github.com/cdzombak/wdproc/releases).

### Build and install locally

```shell
git clone https://github.com/cdzombak/wdproc.git
cd wdproc
make build

cp out/wdproc $INSTALL_DIR
```

## Configuration

Configuration is provided via a YAML file. See [`config.example.yaml`](config.example.yaml).

```yaml
version: 1

paths:
  journal: "~/WriterDeck/Journal"
  tasks: "~/WriterDeck/Tasks"

processing:
  journal_min_age: "4h"
  tasks_min_age: "1h"

integrations:
  day_one:
    enabled: true
    journal: "Default"
    tags: ["via:WriterDeck"]
  
  things:
    enabled: true
    tags: ["via:WriterDeck"]
```

## Usage

```shell
wdproc [OPTIONS]

Options:
  --config PATH     Path to configuration file (default: ./config.yaml)
  --version         Show version information
  --verbose         Enable verbose logging
  --dry-run         Preview actions without executing (default: true)
```

### Basic Usage

1. Configure the tool by editing `config.yaml`
2. Run with `--dry-run=false` to actually process files
3. Files will be processed based on their age and moved to Trash after successful integration

### Prerequisites

- **Day One CLI**: Install via `/usr/local/bin/dayone2` (required for Day One integration)
- **Things 3**: Must be installed for Things integration to work
- **WriterDeck**: Set up your folder structure in `~/WriterDeck/`

## Automation with launchd

For automatic processing, you can use the included launchd plist file `com.dzombak.writerdeckproc.plist` to run wdproc periodically.

**Note**: You will need to adjust the paths in the plist file to include your own username and configuration file location.

```shell
# After making edits, copy the plist to LaunchAgents directory
cp com.dzombak.writerdeckproc.plist ~/Library/LaunchAgents/

# Load the agent (runs every 5 minutes)
launchctl load ~/Library/LaunchAgents/com.dzombak.writerdeckproc.plist

# Start the agent immediately
launchctl start com.dzombak.writerdeckproc

# Check status
launchctl list | grep writerdeckproc
```

## License

MIT License; see [`LICENSE`](LICENSE) in this repo.

## Author

Chris Dzombak ([dzombak.com](https://www.dzombak.com) / [github.com/cdzombak](https://www.github.com/cdzombak)).
