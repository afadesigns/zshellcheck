package reporter

import (
	"encoding/json"
	"io"

	"github.com/afadesigns/zshellcheck/pkg/katas"
)

// SarifReporter reports violations in SARIF format.
type SarifReporter struct {
	writer   io.Writer
	filename string
}

// NewSarifReporter creates a new SarifReporter.
func NewSarifReporter(writer io.Writer, filename string) *SarifReporter {
	return &SarifReporter{
		writer:   writer,
		filename: filename,
	}
}

// Report prints the violations to the writer in SARIF format.
func (r *SarifReporter) Report(violations []katas.Violation) error {
	// Simplified SARIF structure for now
	// Real implementation would be more complex
	type result struct {
		RuleID  string `json:"ruleId"`
		Message string `json:"message"`
		// ... locations, etc.
	}
	type run struct {
		Tool struct {
			Driver struct {
				Name string `json:"name"`
			} `json:"driver"`
		} `json:"tool"`
		Results []result `json:"results"`
	}
	type sarif struct {
		Version string `json:"version"`
		Runs    []run  `json:"runs"`
	}

	results := []result{}
	for _, v := range violations {
		results = append(results, result{
			RuleID:  v.KataID,
			Message: v.Message,
		})
	}

	output := sarif{
		Version: "2.1.0",
		Runs: []run{
			{
				Tool: struct {
					Driver struct {
						Name string `json:"name"`
					} `json:"driver"`
				}{
					Driver: struct {
						Name string `json:"name"`
					}{
						Name: "zshellcheck",
					},
				},
				Results: results,
			},
		},
	}

	encoder := json.NewEncoder(r.writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(output)
}
