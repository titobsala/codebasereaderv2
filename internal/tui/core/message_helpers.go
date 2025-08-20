package core

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// Common message creation helpers following Go's factory function pattern

// NewStatusUpdateMsg creates a status update message
func NewStatusUpdateMsg(message string) tea.Cmd {
	return func() tea.Msg {
		return StatusUpdateMsg{Message: message}
	}
}

// NewErrorMsg creates an error message
func NewErrorMsg(err error) tea.Cmd {
	return func() tea.Msg {
		return ErrorMsg{Error: err}
	}
}

// NewProgressUpdateMsg creates a progress update message
func NewProgressUpdateMsg(current, total int, filePath, message string) tea.Cmd {
	return func() tea.Msg {
		return AnalysisProgressMsg{
			Current:  current,
			Total:    total,
			FilePath: filePath,
			Message:  message,
		}
	}
}

// NewAnalysisStartedMsg creates an analysis started message
func NewAnalysisStartedMsg(path string) tea.Cmd {
	return func() tea.Msg {
		return AnalysisStartedMsg{Path: path}
	}
}

// FormatFileSize formats a file size in bytes to human readable format
func FormatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// FormatProgress formats progress as a percentage string
func FormatProgress(current, total int) string {
	if total == 0 {
		return "0.0%"
	}
	percentage := float64(current) / float64(total) * 100
	return fmt.Sprintf("%.1f%%", percentage)
}
