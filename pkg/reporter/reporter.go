package reporter

import (
	"fmt"
	"io"

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
}

// NewTextReporter creates a new TextReporter.
func NewTextReporter(writer io.Writer, filename string) *TextReporter {
	return &TextReporter{writer: writer, filename: filename}
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
	}
	return nil
}
