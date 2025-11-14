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
	return json.NewEncoder(r.writer).Encode(violations)
}
