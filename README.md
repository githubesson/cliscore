# CliScore (a cli utility for the keyscore api)

A command-line interface for searching data across multiple sources.

## Compilation

### Prerequisites

- Go 1.21 or higher

### Build Instructions

1. Navigate to the project directory:

```bash
cd /path/to/cliscore
```

2. Build the application:

```bash
go build -o cliscore ./cmd/cliscore
```

3. Make the binary executable:

```bash
chmod +x cliscore
```

## Installation

### Add to PATH

#### Method 1: System-wide installation

```bash
sudo mv cliscore /usr/local/bin/
```

#### Method 2: User installation

```bash
mkdir -p ~/bin
mv cliscore ~/bin/
echo 'export PATH="$HOME/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

#### Method 3: Current directory

```bash
export PATH="$(pwd):$PATH"
```

### Verify Installation

```bash
cliscore --help
```

## Configuration

### Environment Variables

- `CLISCORE_BASE_URL`: API base URL (default: https://api.keysco.re)
- `CLISCORE_API_KEY`: Your API key for authentication
- `CLISCORE_RESULTS_DIR`: Directory to save results (default: ~/.keyscore-cli/results)
- `CLISCORE_SAVE_RESULTS`: Enable/disable result saving (true/false)
- `CLISCORE_SPINNER_STYLE`: Spinner style (default, dots, arrows, bounce, simple, none)

### Setup

Run the setup command to configure initial settings:

```bash
cliscore setup
```

## Usage

### Basic Search

```bash
cliscore search <terms>
```

### Advanced Search

```bash
cliscore search -save -operator LOGS test@example.com
```

### Options

- `-source`: Data source to search from (default: xkeyscore)
- `-wildcard`: Enable wildcard search
- `-api-key`: API key for authentication
- `-save`: Save results to file
- `-no-save`: Don't save results to file
- `-results-dir`: Results directory
- `-spinner`: Show loading spinner
- `-quiet`: Quiet mode (no spinner)
- `-operator`: Search operator (AND, LOGS)

### Available Commands

- `search`: Search for terms across different data types
- `count`: Count results for search terms
- `setup`: Configure initial settings
- `config`: Manage configuration
- `machineinfo`: Get machine information
- `download`: Download files or data

## Features

- **Result Counting**: Always displays the number of results found
- **Smart Saving**: Results are only printed to console when saving is disabled
- **Multiple Data Types**: Supports login, password, url, domain, username, ip, hash, phone, uuid
- **Configurable Output**: Choose from different spinner styles or disable entirely
- **File Management**: Automatic safe filename generation and organized storage

## File Locations

- **Config**: `~/.keyscore-cli/config.json`
- **Results**: `~/.keyscore-cli/results/` (or custom directory)
- **Binary**: `/usr/local/bin/cliscore` (or chosen location)
