# Hosts File Parser

[中文版本](README.md)

This is a hosts file parser written in Go, using lexical analysis and syntax analysis techniques to parse and modify hosts files. As a library project, it provides a simple and easy-to-use API for manipulating and managing hosts files.

## Features

- Uses lexical and syntax analysis to parse hosts files
- Supports adding, removing, and finding hosts entries
- Preserves comments and empty lines
- Supports IPv4 and IPv6 address format validation
- Supports mapping multiple domains to a single IP address
- Supports operating hosts entries by hostname
- Supports loading and saving hosts files
- Flexible log system with support for custom logging interfaces

## Installation

```bash
go get github.com/limitcool/hostsparser
```

## Library Usage

### Import Package

```go
import (
    "github.com/limitcool/hostsparser/hosts"
    "github.com/limitcool/hostsparser/logger"
)
```

### Load Hosts File

```go
// Load system hosts file
hostsPath := hosts.GetSystemHostsPath()
hostsFile, err := hosts.LoadHostsFile(hostsPath)
if err != nil {
    // Handle error
}

// Or create a new hosts file
hostsFile := hosts.NewHostsFile()
```

### Modify Hosts Entries

```go
// Set IP for a hostname
err := hostsFile.SetHostIP("example.com", "127.0.0.1")

// Set multiple hostnames to the same IP
err = hostsFile.SetMultipleHostIPs([]string{"example.com", "www.example.com"}, "127.0.0.1")

// Remove hostname mapping
modified, err := hostsFile.RemoveHost("example.com")
```

### Query Hosts Entries

```go
// Get IP for a hostname
ip, err := hostsFile.GetHostIP("example.com")

// Get all hostnames for an IP
domains, err := hostsFile.GetHostsByIP("127.0.0.1")

// Get all IP-domain pairs
pairs := hostsFile.GetAllIPDomainPairs()

// Filter by IP or domain
filteredByIP := hostsFile.FilterIPDomainPairs("127.0.0.1", "")
filteredByDomain := hostsFile.FilterIPDomainPairs("", "example.com")
```

### Save Hosts File

```go
// Save to the original file
err := hostsFile.SaveHostsFile("")

// Save to a new file
err = hostsFile.SaveHostsFile("./hosts.new")
```

### Logging System

This library provides a flexible leveled logging system with support for Debug, Info, Warn, and Error levels.

#### Basic Logging Usage

```go
// Use package-level functions for logging
logger.Debug("This is a debug message")
logger.Info("This is an info message")
logger.Warn("This is a warning message")
logger.Error("This is an error message")

// Use formatting functions
logger.Debugf("Debug: %s", "details")
logger.Infof("Info: %s", "details")
logger.Warnf("Warning: %s", "details")
logger.Errorf("Error: %s", "details")
```

#### Disable Logging

```go
// Completely disable all log output
logger.DisableLogging()
```

#### Custom Logger Interface

You can implement your own logging interface to integrate with high-performance logging libraries like zap, logrus, etc:

```go
// 1. Implement the Logger interface
type MyCustomLogger struct {
    // Custom fields
}

func (l *MyCustomLogger) Debug(args ...interface{}) {
    // Custom implementation
}

func (l *MyCustomLogger) Debugf(format string, args ...interface{}) {
    // Custom implementation
}

// Implement other methods...

// 2. Set as the global logger
logger.SetLogger(&MyCustomLogger{})
```

#### Zap Logger Adapter Example

```go
import (
    "github.com/limitcool/hostsparser/logger"
    "go.uber.org/zap"
)

// In the main program
func main() {
    // Create a zap logger instance
    zapLogger, _ := zap.NewProduction()
    defer zapLogger.Sync()

    // Use zap's Sugar interface to adapt to our Logger interface
    logger.SetLogger(zapLogger.Sugar())

    // Now all logs will be recorded through zap
}
```

## Project Structure

- `lexer/` - Lexical analyzer that breaks down hosts file content into tokens
- `parser/` - Syntax analyzer that parses tokens and builds hosts entry structures
- `hosts/` - Core functionality for hosts file operations
- `logger/` - Leveled logging system with support for custom logging interfaces

## Implementation Details

### Lexical and Syntax Analysis

This project uses classic compiler frontend techniques to parse hosts files:

1. **Lexical Analysis**: Breaks down input text into meaningful tokens, such as IP addresses, domains, comments, etc.

2. **Syntax Analysis**: Analyzes token sequences and builds structured hosts entries.

This approach makes the parsing process more robust, capable of handling various formats of hosts files, including comments, empty lines, and malformed lines.

### Hosts File Format

Standard hosts file format:

```
# This is a comment line
127.0.0.1 localhost
::1 localhost ipv6-localhost ipv6-loopback
192.168.1.1 example.com www.example.com # This is a trailing comment
```

## License

MIT

## Contributing

Pull Requests and Issues are welcome.
