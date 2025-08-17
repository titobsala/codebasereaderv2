# TUI Stability and Usability Overhaul

This ticket consolidates several critical bugs and key enhancements identified during recent testing. The issues range from a fatal runtime panic during analysis to severe UI/UX problems with navigation and item counting.

**Labels:** `bug`, `enhancement`, `TUI`, `panic`, `usability`

-----

## 1\. Bug: Goroutine Panic on Analysis - "close of closed channel"

The application crashes with a fatal panic when running the codebase analysis. The stack trace points to a race condition in the worker pool's stop mechanism.

**Severity:** Blocker 🛑

### Steps to Reproduce

1.  Launch the application in a valid project directory.
2.  Select one or more files/folders for analysis.
3.  Start the analysis process.

### Expected Behavior

The analysis process should complete successfully and display the results, or be cancelled cleanly by the user without crashing.

### Actual Behavior

The application terminates unexpectedly, printing the following panic to the console. This suggests that the `(*WorkerPool).Stop` function is attempting to close a channel that has already been closed, which is a common concurrency issue.

```go
Caught panic:

close of closed channel

Restoring terminal...

goroutine 54 [running]:
runtime/debug.Stack()
        /usr/local/go/src/runtime/debug/stack.go:26 +0x5e
runtime/debug.PrintStack()
        /usr/local/go/src/runtime/debug/stack.go:18 +0x13
github.com/charmbracelet/bubbletea.(*Program).recoverFromPanic(0xc0000ef180)
        /home/tito-sala/go/pkg/mod/github.com/charmbracelet/bubbletea@v1.3.4/tea.go:758 +0x8b
panic({0x5c3ea0?, 0x654e20?})
        /usr/local/go/src/runtime/panic.go:783 +0x132
github.com/tito-sala/codebasereaderv2/internal/engine.(*WorkerPool).Stop(0xc0000ed180)
        /home/tito-sala/Code/Personal/codebasereaderv2/internal/engine/engine.go:253 +0x6c
github.com/tito-sala/codebasereaderv2/internal/engine.(*Engine).AnalyzeDirectoryWithProgress(0xc0000bfbf0, {0xc0003aac00, 0x2e}, 0x0)
        /home/tito-sala/Code/Personal/codebasereaderv2/internal/engine/engine.go:117 +0x8ba
github.com/tito-sala/codebasereaderv2/internal/engine.(*Engine).AnalyzeDirectoryWithEnhancedMetricsAndProgress(0xc0000bfbf0, {0xc0003aac00, 0x2e}, 0xc0000a22e0?)
        /home/tito-sala/Code/Personal/codebasereaderv2/internal/engine/engine.go:368 +0x3a
github.com/tito-sala/codebasereaderv2/internal/engine.(*Engine).AnalyzeDirectoryWithEnhancedMetrics(0xc000602c40?, {0xc0003aac00?, 0xc0000309a0?})
        /home/tito-sala/Code/Personal/codebasereaderv2/internal/engine/engine.go:362 +0x1a
github.com/tito-sala/codebasereaderv2/internal/tui.(*MainModel).startAnalysis.(*MainModel).performAnalysis.func2()
        /home/tito-sala/Code/Personal/codebasereaderv2/internal/tui/model.go:724 +0x3b
github.com/charmbracelet/bubbletea.(*Program).handleCommands.func1.1()
        /home/tito-sala/go/pkg/mod/github.com/charmbracelet/bubbletea@v1.3.4/tea.go:352 +0x7d
created by github.com/charmbracelet/bubbletea.(*Program).handleCommands.func1 in goroutine 30
        /home/tito-sala/go/pkg/mod/github.com/charmbracelet/bubbletea@v1.3.4/tea.go:346 +0x131
2025/08/17 21:58:32 program was killed: context canceled
```

-----

## 2\. Bug: TUI Navigation is Unintuitive and State Becomes Inconsistent

The file navigator's behavior with left/right arrow keys is broken. The left arrow key incorrectly changes the working directory (`cd ..`) instead of collapsing the selected folder in the tree view. This leads to a desynchronization between the UI display and the application's internal state.

**Severity:** High  крити

### Steps to Reproduce

1.  Start the application at a path, e.g., `/home/tito-sala/Code/Personal/codebasereaderv2`.
2.  Navigate down to a subdirectory (e.g., `cmd`).
3.  Press the right arrow key to expand the `cmd` directory. The view correctly updates.
4.  Press the left arrow key.
      - **Observed:** The view changes to the parent directory `/home/tito-sala/Code/Personal`. The header path may or may not update correctly.
      - **Expected:** The `cmd` directory should simply collapse in the tree view, with the focus remaining on it.
5.  Press the right arrow key again.
      - **Observed:** The view now displays the contents of `codebasereaderv2` again, but the header path likely still shows `/home/tito-sala/Code/Personal`. The state is now inconsistent.
6.  Pressing the left arrow key again will navigate to `/home/tito-sala/Code`, further demonstrating the incorrect `cd ..` behavior.

-----

## 3\. Bug: TUI Footer Item Counter is Static

The item counter in the footer (e.g., `1/12`) does not update when the user scrolls through the list of files and directories.

**Severity:** Low  cosmetic

### Steps to Reproduce

1.  Navigate to a directory containing more items than can be displayed on a single screen.
2.  Observe the initial item count in the footer.
3.  Use the up/down arrow keys to scroll through the list.

### Expected Behavior

The counter should dynamically update to reflect the index of the currently highlighted item (e.g., `2/12`, `3/12`, etc.).

### Actual Behavior

The counter remains frozen on its initial value, providing no feedback on the user's position in the list.

-----

## 4\. Enhancement: Default Analysis to Current Directory

To improve user workflow, the analysis feature should be enhanced. If a user triggers an analysis without explicitly toggling any files or folders, the application should default to analyzing the entire current directory being viewed.

### User Story

As a user, when I want to analyze an entire project, I want to be able to simply press the "analyze" key without any selections, so that the tool intelligently analyzes the current working directory, saving me time.

### Acceptance Criteria

  - **Given** a user is viewing a directory in the TUI
  - **And** no files or folders are toggled/selected
  - **When** the user initiates an analysis
  - **Then** the analysis should run on the root of the current directory.
  - **And** if one or more items *are* toggled, the analysis should proceed only on those selected items (maintaining current behavior).