# TUI/Engine Polish and Scalability Improvements

This ticket addresses follow-up issues identified after initial fixes. The focus is on refining navigation, fixing a critical analysis-engine scaling problem, and correcting a UI rendering bug in the file tree.

**Labels:** `bug`, `enhancement`, `TUI`, `engine`, `scalability`

-----

## 1\. Bug: TUI - Incorrect Indentation in File Tree Display

The file tree renderer incorrectly indents root-level files, making them appear as children of the last listed root-level directory, which creates a confusing and inaccurate representation of the project structure.

**Severity:** Medium 🎨

### Description

Files that exist in the root of the project directory are being displayed visually as if they are nested inside another directory. For instance, files like `README.md` and `go.mod` are shown at the same indentation level as the contents of the `internal` folder, not at the top level alongside `bin`, `cmd`, and `internal`.

### Current Display (Incorrect)

As seen in the image, the root files are misaligned:

### Expected Behavior

All files and folders in the root directory should be aligned at the same, outermost indentation level. The hierarchy should be clear and accurate.

```
/
├── 📁 bin
├── 📁 cmd
├── 📁 internal
├── 📝 CLAUDE.md
├── 📝 ENHANCED_METRICS_SUMMARY.md
├── ...
└── ❓ go.mod
```

-----

## 2\. Bug: Engine Fails on Full Directory Analysis - "job queue is full"

When attempting to analyze a full directory (the default behavior with no items toggled), the analysis engine fails immediately. The error indicates that the concurrent job processing queue is not large enough to handle all the files in the target directory.

**Severity:** High

### Steps to Reproduce

1.  Launch the application in a moderately sized project directory.
2.  Ensure no specific files or folders are toggled.
3.  Trigger the analysis.

### Expected Behavior

The analysis should begin processing all valid files within the current directory and its subdirectories.

### Actual Behavior

The process terminates with an error, preventing the analysis from running. The error message is:
`Error: failed to submit job for /home/tito-sala/Code/Personal/codebasereaderv2/cmd/tui/main.go: job queue is full`

### Suggested Action

The worker pool's job queue size needs to be re-evaluated. Consider making it larger, dynamic, or implementing a backpressure mechanism to prevent it from being overwhelmed.

### Future Considerations

To make analysis more efficient, we should implement a file exclusion mechanism, similar to `.gitignore`. This would allow the tool to skip directories like `node_modules`, `__pycache__`, vendor folders, and build artifacts, reducing the number of jobs submitted to the queue.

-----

## 3\. Enhancement: Implement Confirmed "Go Up" Directory Navigation

The previous "back" behavior (left arrow) was buggy. This proposal refines it into an intentional feature: allowing the user to navigate to the parent directory (`cd ..`) from the TUI, but only after a confirmation prompt to prevent accidental navigation.

**Severity:** Medium ✨

### User Story

As a user, when I am viewing a directory's file list, I want to be able to press the left arrow key to navigate up to the parent directory, but only after confirming my intent, so I don't leave my current context by mistake.

### Acceptance Criteria

1.  When the user presses the **left arrow key** on the **first item** in the file tree (e.g., a "../" entry or just by being at the top), a confirmation prompt should appear.
2.  The prompt should be clear, for example: "**Navigate to parent directory? (y/N)**".
3.  If the user confirms ('y'), the application view should change to the parent directory.
4.  If the user cancels ('n' or Esc), the prompt should disappear, and the view should remain unchanged.
5.  Pressing the left arrow on any other item (e.g., an expanded folder) should perform the "collapse" action as expected, without a prompt.
6.  The UI for the prompt can be implemented using components from the **Charm Bracelet** libraries (`bubbletea`, `lipgloss`) for a consistent look and feel.