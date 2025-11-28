package reporter

import (
	"encoding/json"
	"io"

	"github.com/afadesigns/zshellcheck/pkg/katas"
)

// JSONReporter is a reporter that writes JSON to an io.Writer.
type JSONReporter struct {
	writer io.Writer
}

// NewJSONReporter creates a new JSONReporter.
func NewJSONReporter(writer io.Writer) *JSONReporter {
	return &JSONReporter{writer: writer}
}

// Report prints the violations to the writer as a JSON array.
func (r *JSONReporter) Report(violations []katas.Violation) error {
	type jsonViolation struct {
		katas.Violation
		Title string `json:"Title"`
		Level string `json:"Level"`
	}

	var output []jsonViolation
	for _, v := range violations {
		kata, ok := katas.Registry.GetKata(v.KataID)
		title := ""
		if ok {
			title = kata.Title
		}
		output = append(output, jsonViolation{
			Violation: v,
			Title:     title,
			Level:     "warning", // Default level for now
		})
	}

	enc := json.NewEncoder(r.writer)
	enc.SetIndent("", "  ")
	return enc.Encode(output)
}
