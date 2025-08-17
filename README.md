# CodebaseReader v2 🔍

A powerful, interactive Terminal User Interface (TUI) for codebase analysis and exploration. Built with Go and designed for developers who need to quickly understand and navigate unfamiliar codebases.

![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)
![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)

## ✨ Features

### 🖥️ Interactive Terminal Interface

- **Beautiful TUI** built with [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **File tree navigation** with expand/collapse functionality
- **Multiple view modes** for different types of analysis
- **Real-time progress indicators** during analysis
- **Keyboard shortcuts** for efficient navigation

### ⚡ High-Performance Analysis

- **Concurrent processing** using Go's goroutines and worker pools
- **Multi-language support** with pluggable parser architecture
- **Smart file filtering** with .gitignore support
- **Memory-efficient** processing of large codebases

### 📊 Comprehensive Code Metrics

- **Project-level statistics**: Total files, lines of code, complexity scores
- **Language breakdown**: Per-language metrics and distribution
- **Function and class analysis**: Detailed code structure insights
- **File-level details**: Individual file metrics and analysis

### 🎯 Currently Supported Languages

- **Go** (.go) - Full support with AST parsing
- **Python** (.py) - _Coming soon_
- **JavaScript/TypeScript** - _Planned_
- **Java** - _Planned_

## 🚀 Quick Start

### Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/your-username/codebasereaderv2.git
   cd codebasereaderv2
   ```

2. **Build the application:**

   ```bash
   go build -o codebase-analyzer ./cmd/tui
   ```

3. **Run the analyzer:**
   ```bash
   ./codebase-analyzer
   ```

### Usage

Once the TUI is running, you can:

1. **Navigate** through your file system using the file tree
2. **Select a directory** and press `a` to analyze it
3. **View results** in multiple formats (overview, detailed metrics)
4. **Switch between views** using keyboard shortcuts

## ⌨️ Keyboard Shortcuts

### Navigation

| Key              | Action                          |
| ---------------- | ------------------------------- |
| `↑/↓` or `k/j`   | Move up/down in file tree       |
| `→/l` or `Enter` | Expand directory or select file |
| `←/h`            | Collapse directory or go back   |
| `Tab`            | Switch between views            |
| `Esc`            | Return to file tree view        |

### Analysis

| Key | Action                       |
| --- | ---------------------------- |
| `a` | Analyze selected directory   |
| `r` | Refresh file tree            |
| `m` | Toggle detailed metrics view |
| `s` | Toggle summary view          |

### General

| Key             | Action           |
| --------------- | ---------------- |
| `?`             | Show/hide help   |
| `q` or `Ctrl+C` | Quit application |

## 📈 Analysis Output

### Overview Mode

```
📊 Codebase Analysis Results
========================================

📁 Root Path: /path/to/project
📄 Total Files: 42
📝 Total Lines: 3,847
🕒 Generated: 2024-08-16 15:30:45

🌐 Languages:
  Go:
    Files: 38
    Lines: 3,521
    Functions: 156
    Classes: 23
    Complexity: 89
```

### Detailed Metrics Mode

```
📈 Detailed Metrics
========================================

📊 Overall Statistics:
  Total Files: 42
  Total Lines: 3,847
  Total Functions: 156
  Total Classes: 23

🔧 Go Metrics:
  Files: 38 (90.5%)
  Lines: 3,521 (91.5%)
  Functions: 156
  Classes: 23
  Avg Lines/Function: 22.6
  Avg Complexity/Function: 2.3
```

## 🏗️ Architecture

The project follows a modular architecture with clear separation of concerns:

```
├── cmd/
│   └── tui/           # TUI application entry point
├── internal/
│   ├── engine/        # Analysis engine and worker pools
│   ├── parser/        # Language parsers (Go, Python, etc.)
│   ├── tui/          # Terminal UI components
│   └── config/       # Configuration management
```

### Key Components

- **Analysis Engine**: Concurrent file processing with worker pools
- **Parser Registry**: Pluggable system for adding new language parsers
- **TUI Framework**: Interactive terminal interface built with Bubble Tea
- **File Walker**: Efficient directory traversal with filtering

## 🛠️ Development

### Prerequisites

- Go 1.21 or higher
- Terminal with color support

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/engine/...
```

### Building from Source

```bash
# Build for current platform
go build -o codebase-analyzer ./cmd/tui

# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o codebase-analyzer-linux ./cmd/tui
GOOS=darwin GOARCH=amd64 go build -o codebase-analyzer-macos ./cmd/tui
GOOS=windows GOARCH=amd64 go build -o codebase-analyzer-windows.exe ./cmd/tui
```

## 🗺️ Roadmap

### ✅ Completed

- [x] Interactive TUI with file tree navigation
- [x] Go language parser with AST analysis
- [x] Concurrent analysis engine
- [x] Real-time progress reporting
- [x] Multiple view modes for results
- [x] Comprehensive code metrics

### 🚧 In Progress

- [ ] Python language parser
- [ ] AI-powered code summaries
- [ ] Export functionality (JSON, Mermaid diagrams)

### 📋 Planned

- [ ] JavaScript/TypeScript support
- [ ] Java language parser
- [ ] Command-line interface (headless mode)
- [ ] Configuration file support
- [ ] Plugin system for custom parsers
- [ ] Performance optimizations and caching

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

This project is built with amazing open-source libraries:

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components

## 📞 Support

If you have any questions or run into issues, please:

1. Check the [Issues](https://github.com/your-username/codebasereaderv2/issues) page
2. Create a new issue with detailed information
3. Include your Go version and operating system

---

**Made with ❤️ and Go**
