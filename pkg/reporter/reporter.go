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
	writer io.Writer
}

// NewTextReporter creates a new TextReporter.
func NewTextReporter(writer io.Writer) *TextReporter {
	return &TextReporter{writer: writer}
}

// Report prints the violations to the writer.
func (r *TextReporter) Report(violations []katas.Violation) error {
	for _, v := range violations {
		kata, ok := katas.GetKata(v.KataID)
		if !ok {
			return fmt.Errorf("kata with ID %s not found", v.KataID)
		}
		_, err := fmt.Fprintf(r.writer, "%s: %s (%s)\n", v.KataID, v.Message, kata.Title)
		if err != nil {
			return err
		}
	}
	return nil
}
