# Obsidian Sync

Go-based file watcher that monitors Obsidian vault changes and sends events to a
cloud API. Built for creating AI-powered knowledge systems with real-time file
synchronization.

## Features

- 🔍 **Recursive File Monitoring**: Watches entire Obsidian vault including
  subdirectories
- ⚡ **Event Debouncing**: Handles rapid file system events from text editors
  intelligently
- 🎯 **Smart Event Detection**: Distinguishes between create, modify, and delete
  operations
- 📝 **Atomic Save Handling**: Properly handles editor save patterns (rename →
  create)
- 📊 **Structured Logging**: Comprehensive logging with rotation and proper
  caller information
- ⚙️ **Environment Configuration**: `.env` file support with validation
- 🔐 **API Integration**: Ready to send authenticated requests to cloud APIs

## Prerequisites

- Go 1.21 or higher
- An Obsidian vault
- Access to a cloud API endpoint (see
  [obsidian-infrastructure](../obsidian-infrastructure))

## Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/obsidian-sync.git
cd obsidian-sync
```

2. Install dependencies:

```bash
go mod tidy
```

3. Copy the environment template:

```bash
cp .env.example .env
```

4. Configure your environment variables in `.env`:

```bash
# Required
VAULT_PATH=/path/to/your/obsidian/vault
API_ENDPOINT=https://your-api-gateway-url.com/events
API_KEY=your-api-key

# Optional
LOG_LEVEL=info
LOG_FILE=logs/obsidian-sync.log
```

## Usage

### Run the File Watcher

```bash
# Development
go run cmd/sync/main.go

# Build and run
go build -o obsidian-sync cmd/sync/main.go
./obsidian-sync
```

### Running with Docker

```bash
# Build image
docker build -t obsidian-sync .

# Run container
docker run -v /path/to/vault:/vault -v /path/to/logs:/app/logs obsidian-sync
```

## Configuration

### Environment Variables

| Variable       | Description                              | Default                  | Required |
| -------------- | ---------------------------------------- | ------------------------ | -------- |
| `VAULT_PATH`   | Path to your Obsidian vault              | -                        | Yes      |
| `API_ENDPOINT` | Cloud API endpoint URL                   | -                        | Yes      |
| `API_KEY`      | API authentication key                   | -                        | Yes      |
| `LOG_LEVEL`    | Logging level (debug, info, warn, error) | `info`                   | No       |
| `LOG_FILE`     | Path to log file                         | `logs/obsidian-sync.log` | No       |

### Logging Configuration

The application uses structured logging with automatic rotation:

- **Log Rotation**: 100MB per file, keeps 10 old files
- **Retention**: 30 days
- **Compression**: Old logs are compressed
- **Output**: Both console and file logging

## Project Structure

```
obsidian-sync/
├── cmd/
│   └── sync/
│       └── main.go          # Application entry point
├── internal/
│   ├── watcher/
│   │   └── watcher.go       # File monitoring logic
│   ├── config/
│   │   └── config.go        # Configuration management
│   ├── logger/
│   │   └── logger.go        # Logging setup
│   └── client/
│       └── api.go           # HTTP client (coming soon)
├── pkg/
│   └── models/
│       └── file.go          # Shared data structures
├── .env.example             # Environment template
├── go.mod                   # Go module file
└── README.md
```

## Event Format

The watcher sends JSON events in this format:

```json
{
  "event_type": "file_modified",
  "file_path": "/Users/username/vault/daily-notes/2025-06-08.md",
  "vault_path": "/Users/username/vault",
  "timestamp": "2025-06-08T14:30:00Z",
  "file_size": 1024,
  "checksum": "abc123def456"
}
```

### Event Types

- `file_created`: New file added to vault
- `file_modified`: Existing file changed
- `file_deleted`: File removed from vault

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/watcher
```

### Debugging

Set log level to `debug` for detailed information:

```bash
LOG_LEVEL=debug go run cmd/sync/main.go
```

## Roadmap

- [ ] HTTP client implementation for API requests
- [ ] Initial vault synchronization
- [ ] Retry logic and error handling
- [ ] File content diffing for incremental updates
- [ ] Metadata extraction (tags, links, backlinks)
- [ ] Performance optimizations for large vaults

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file
for details.

## Related Projects

- [obsidian-infrastructure](../obsidian-infrastructure) - AWS infrastructure for
  the complete RAG system
- [obsidian-pipeline](../obsidian-pipeline) - Lambda functions for processing
  and embedding files
