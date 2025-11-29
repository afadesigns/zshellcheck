package reporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/config"
	"github.com/afadesigns/zshellcheck/pkg/katas"
)

// Reporter defines the interface for reporting violations.
type Reporter interface {
	Report(violations []katas.Violation) error
}

// TextReporter is a simple reporter that writes plain text to an io.Writer.
type TextReporter struct {
	writer   io.Writer
	filename string
	lines    []string
	config   config.Config
}

// NewTextReporter creates a new TextReporter.
func NewTextReporter(writer io.Writer, filename, source string, config config.Config) *TextReporter {
	return &TextReporter{
		writer:   writer,
		filename: filename,
		lines:    strings.Split(source, "\n"),
		config:   config,
	}
}

func (r *TextReporter) getColor(code string) string {
	if r.config.NoColor {
		return ""
	}
	return code
}

// Report prints the violations to the writer.
func (r *TextReporter) Report(violations []katas.Violation) error {
	for _, v := range violations {
		// Severity Color
		color := ""
		switch v.Level {
		case katas.Error:
			color = r.getColor("\033[31m") // Red
		case katas.Warning:
			color = r.getColor("\033[33m") // Yellow
		case katas.Info:
			color = r.getColor("\033[34m") // Blue
		}
		reset := r.getColor("\033[0m")
		bold := r.getColor("\033[1m")

		// Location: filename:line:col
		fmt.Fprintf(r.writer, "%s:%d:%d: ", r.filename, v.Line, v.Column)

		// Severity, ID, Message
		// Example: Error: [ZC1001] Some message
		fmt.Fprintf(r.writer, "%s%s%s: [%s] %s\n", color, v.Level, reset, v.KataID, v.Message)

		// Code snippet
		if v.Line > 0 && v.Line <= len(r.lines) {
			lineContent := r.lines[v.Line-1]
			// Replace tabs with spaces for correct alignment of caret (simple approach)
			// Or keep it simple for now.
			fmt.Fprintf(r.writer, "  %s\n", lineContent)

			// Caret
			padding := v.Column - 1
			if padding < 0 {
				padding = 0
			}
			// Use a simple space padding. Note: this might be slightly off if tabs are present,
			// but it's a standard starting point.
			fmt.Fprintf(r.writer, "  %s%s^%s\n", strings.Repeat(" ", padding), bold, reset)
		}
		fmt.Fprintln(r.writer)
	}
	return nil
}