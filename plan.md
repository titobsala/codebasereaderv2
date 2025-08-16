Of course. This is an excellent project, and rewriting it in Go with a more extensible design is a great professional development exercise.

Here is a Project Design Document (PDD) that outlines the foundation for building your new, multi-language `CodebaseReader` using Go and the Charm libraries.

-----

# Project Design Document: CodebaseReader v2

## 1\. Introduction

### 1.1 Project Goal

To re-engineer the `CodebaseReader` tool in the Go programming language, transforming it from a single-language (Python) script into a high-performance, extensible, multi-language analysis tool with an interactive Terminal User Interface (TUI).

### 1.2 Background

The original `CodebaseReader` was effective for analyzing Python projects but was limited by its core dependency on Python's `ast` library. This rewrite is motivated by the desire to:

1.  Support a variety of programming languages.
2.  Improve performance through concurrency.
3.  Provide a superior user experience with a modern TUI.
4.  Distribute the application as a single, portable binary.

### 1.3 Key Features

  * Interactive TUI for directory navigation and viewing analysis results.
  * Concurrent analysis of files for high performance.
  * **Pluggable architecture** for easily adding new language parsers.
  * Generation of code metrics (line count, function/class count, complexity).
  * Visualization of project structure and code relationships (e.g., Mermaid diagrams).
  * Integration with AI services for automated code summarization.

-----

## 2\. System Architecture & Design

### 2.1 Core Philosophy: The Parser Plugin Pattern

The central challenge is analyzing different languages. We will solve this by designing a **pluggable architecture**. The core application will be language-agnostic. It will treat language analysis as a service provided by "parser plugins."

Each supported language will have its own parser that conforms to a standard **Go interface**. This makes the system incredibly extensible: to support a new language, we simply create a new struct that satisfies the `Parser` interface.

### 2.2 High-Level Diagram

```
+--------------------------+
|      TUI (Bubble Tea)    |
| (User Input, Display)    |
+-------------+------------+
              |
+-------------v------------+
|      Core Engine         |
| (File Walker, Conductor) |
+-------------+------------+
              | Dispatches based on file type (*.py, *.js, *.go)
+-------------v------------+      +-----------------+      +--------------------+
|   Parser Interface       |----->|  Python Parser  |----->| go-python/gpython  |
| - GetFunctions()         |      +-----------------+      +--------------------+
| - GetClasses()           |
| - GetMetrics()           |      +-----------------+      +--------------------+
| ...                      |----->| JavaScript Parser|----->| third-party JS lib |
+--------------------------+      +-----------------+      +--------------------+
                                  |      ...        |
```

### 2.3 Component Breakdown

  * **TUI Frontend (`bubbletea`, `lipgloss`):**

      * This is the user-facing part of the application.
      * It will consist of components like a file/directory tree navigator, a main content pane to display analysis results, a status bar for feedback, and input fields for configuration (e.g., API keys).

  * **Core Engine:**

      * This is the central orchestrator.
      * It's responsible for walking the specified directory, discovering files, and managing a pool of goroutines for concurrent processing.
      * It maintains a **parser registry** (e.g., a `map[string]Parser`) that maps file extensions (`.py`, `.go`) to the appropriate parser implementation.

  * **Parser Interface:**

      * This is the cornerstone of our extensible design. We'll define a Go `interface` that every language parser must implement.

    <!-- end list -->

    ```go
    // a simplified example
    type AnalysisResult struct {
        FunctionName string
        LineCount    int
        Complexity   int
    }

    type Parser interface {
        // Parse takes file content and returns structured analysis data.
        Parse(fileContent []byte) ([]AnalysisResult, error)
    }
    ```

  * **AI Summarizer:**

      * A separate module responsible for taking the collected code snippets or analysis data, formatting it into a prompt, and sending it to an AI API (like Anthropic or an OpenAI-compatible endpoint).

-----

## 3\. Functional Requirements (FR)

  * **FR1:** The application MUST start by displaying an interactive file tree of the current directory.
  * **FR2:** The user MUST be able to navigate the file tree using the keyboard.
  * **FR3:** Upon selecting a directory, the application MUST concurrently scan for supported file types (`.go`, `.py`, etc.).
  * **FR4:** For each supported file, the Core Engine MUST delegate analysis to the corresponding registered parser.
  * **FR5:** The TUI MUST display aggregated metrics for the selected directory (e.g., total files, total lines of code, list of all classes/functions).
  * **FR6:** The application MUST provide an option to generate a technical summary of the codebase using an external AI API.
  * **FR7:** The application MUST be configurable via a configuration file or command-line flags to set the AI provider and API key.

-----

## 4\. Non-Functional Requirements (NFR)

  * **NFR1 (Performance):** Analysis of large codebases should be fast. File I/O and parsing must be done concurrently.
  * **NFR2 (Extensibility):** Adding support for a new language (e.g., Ruby) should NOT require changes to the Core Engine or TUI—only the creation of a new file (`ruby_parser.go`) that implements the `Parser` interface.
  * **NFR3 (Usability):** The TUI must be responsive and provide clear feedback to the user (e.g., loading spinners, progress bars).
  * **NFR4 (Portability):** The project MUST be buildable into a single, static binary with no external runtime dependencies.

-----

## 5\. Technology Stack

  * **Language:** Go
  * **TUI Framework:** Charm Bubble Tea & Lipgloss
  * **Concurrency:** Go Goroutines & Channels
  * **Parsing Libraries:** To be determined for each language (e.g., `go/parser` for Go, a third-party library for Python).

-----

## 6\. Development Roadmap

  * **Phase 1: The Core & First Parser**

    1.  Define the final `Parser` interface.
    2.  Build the non-TUI Core Engine that can walk a directory.
    3.  Implement the **Go parser** first. Using the native `go/parser` package will be the easiest, providing a quick win and a solid test case.
    4.  Make the core engine produce a simple JSON output of its findings.

  * **Phase 2: Build the TUI**

    1.  Design the basic TUI layout with `bubbletea`.
    2.  Integrate the Core Engine from Phase 1. Make the TUI call the engine and display the JSON results in a formatted way.
    3.  Add interactive elements like the file tree.

  * **Phase 3: Expand Language Support & AI**

    1.  Research and integrate a third-party Go library to parse **Python code**.
    2.  Implement the Python parser struct that satisfies the `Parser` interface.
    3.  Build the AI Summarizer module and connect it to the TUI.

  * **Phase 4: Refinement & Distribution**

    1.  Add more output formats (e.g., Mermaid diagrams).
    2.  Refine performance and add robust error handling.
    3.  Set up a build pipeline (e.g., GitHub Actions) to compile and release the single binary for different operating systems.
