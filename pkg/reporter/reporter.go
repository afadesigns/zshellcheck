package reporter

import (
	"fmt"
	"io"
	"strings"

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
}

// NewTextReporter creates a new TextReporter.
func NewTextReporter(writer io.Writer, filename, source string) *TextReporter {
	return &TextReporter{
		writer:   writer,
		filename: filename,
		lines:    strings.Split(source, "\n"),
	}
}

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

// Report prints the violations to the writer.
func (r *TextReporter) Report(violations []katas.Violation) error {
	for _, v := range violations {
		kata, ok := katas.Registry.GetKata(v.KataID)
		if !ok {
			return fmt.Errorf("kata with ID %s not found", v.KataID)
		}

		// Format: file:line:col: [ID] Title: Message
		// Example: demo.zsh:10:5: [ZC1001] Array Access: Invalid array access...

		fmt.Fprintf(r.writer, "%s%s:%d:%d:%s %s[%s]%s %s%s:%s %s\n",
			colorBold, r.filename, v.Line, v.Column, colorReset,
			colorRed, v.KataID, colorReset,
			colorCyan, kata.Title, colorReset,
			v.Message,
		)

		// Print source line context with gutter
		if v.Line > 0 && v.Line <= len(r.lines) {
			line := r.lines[v.Line-1]
			lineNumStr := fmt.Sprintf("%d", v.Line)
			
			// Gutter padding based on line number length
			// For simple alignment, we can assume line numbers won't differ wildly in width for adjacent errors,
			// but let's just print it directly.
			
			fmt.Fprintf(r.writer, "  %s%s |%s %s\n", colorCyan, lineNumStr, colorReset, line)
			
			// Calculate pointer padding
			pad := ""
			for i := 0; i < v.Column-1 && i < len(line); i++ {
				if line[i] == '\t' {
					pad += "\t"
				} else {
					pad += " "
				}
			}
			
			gutterSpace := strings.Repeat(" ", len(lineNumStr))
			fmt.Fprintf(r.writer, "  %s %s|%s %s%s^%s\n", gutterSpace, colorCyan, colorReset, pad, colorYellow, colorReset)
		}
		fmt.Fprintln(r.writer) // Add blank line between violations
	}
	return nil
}
