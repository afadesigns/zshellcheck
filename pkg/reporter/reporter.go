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

		// Format: file:line:col: [ID] Message (Title)
		// Example: demo.zsh:10:5: [ZC1001] Invalid array access (Use ${var}...)

		fmt.Fprintf(r.writer, "%s%s:%d:%d:%s %s[%s]%s %s %s(%s)%s\n",
			colorBold, r.filename, v.Line, v.Column, colorReset,
			colorRed, v.KataID, colorReset,
			v.Message,
			colorCyan, kata.Title, colorReset,
		)

		// Print source line context
		if v.Line > 0 && v.Line <= len(r.lines) {
			line := r.lines[v.Line-1]
			pad := ""
			for i := 0; i < v.Column-1 && i < len(line); i++ {
				if line[i] == '\t' {
					pad += "\t"
				} else {
					pad += " "
				}
			}
			fmt.Fprintf(r.writer, "  %s\n", line)
			fmt.Fprintf(r.writer, "  %s%s^%s\n", pad, colorYellow, colorReset)
		}
	}
	return nil
}
