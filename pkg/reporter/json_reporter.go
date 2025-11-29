package reporter

import (
	"encoding/json"
	"io"

	"github.com/afadesigns/zshellcheck/pkg/katas"
)

// JSONReporter reports violations in JSON format.
type JSONReporter struct {
	writer io.Writer
}

// NewJSONReporter creates a new JSONReporter.
func NewJSONReporter(writer io.Writer) *JSONReporter {
	return &JSONReporter{writer: writer}
}

// Report prints the violations to the writer in JSON format.
func (r *JSONReporter) Report(violations []katas.Violation) error {
	encoder := json.NewEncoder(r.writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(violations)
}
