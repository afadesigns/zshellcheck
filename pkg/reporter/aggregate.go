// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package reporter

import (
	"encoding/json"
	"io"

	"github.com/afadesigns/zshellcheck/pkg/katas"
)

// FileViolations pairs a scanned file with the violations found in it.
// The machine-readable reporters aggregate these across every scanned
// file into one document so the output is a single valid JSON / SARIF
// value, with each finding attributed to its file.
type FileViolations struct {
	Filename   string
	Violations []katas.Violation
}

type jsonFinding struct {
	File    string         `json:"File"`
	KataID  string         `json:"KataID"`
	Message string         `json:"Message"`
	Line    int            `json:"Line"`
	Column  int            `json:"Column"`
	Level   katas.Severity `json:"Level"`
}

// ReportJSON writes every finding across all files as one JSON array.
// Each element keeps the original single-file fields and adds `File`, so
// existing single-file consumers are unaffected and multi-file output is
// valid and attributed.
func ReportJSON(w io.Writer, files []FileViolations) error {
	findings := []jsonFinding{}
	for _, f := range files {
		for _, v := range f.Violations {
			findings = append(findings, jsonFinding{
				File:    f.Filename,
				KataID:  v.KataID,
				Message: v.Message,
				Line:    v.Line,
				Column:  v.Column,
				Level:   v.Level,
			})
		}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(findings)
}

// SARIF 2.1.0 document shape, trimmed to the fields ZShellCheck emits.
type sarifDoc struct {
	Schema  string     `json:"$schema"`
	Version string     `json:"version"`
	Runs    []sarifRun `json:"runs"`
}

type sarifRun struct {
	Tool    sarifTool     `json:"tool"`
	Results []sarifResult `json:"results"`
}

type sarifTool struct {
	Driver sarifDriver `json:"driver"`
}

type sarifDriver struct {
	Name           string `json:"name"`
	InformationURI string `json:"informationUri"`
	Version        string `json:"version"`
}

type sarifResult struct {
	RuleID    string          `json:"ruleId"`
	Level     string          `json:"level"`
	Message   sarifMessage    `json:"message"`
	Locations []sarifLocation `json:"locations"`
}

type sarifMessage struct {
	Text string `json:"text"`
}

type sarifLocation struct {
	PhysicalLocation sarifPhysical `json:"physicalLocation"`
}

type sarifPhysical struct {
	ArtifactLocation sarifArtifact `json:"artifactLocation"`
	Region           sarifRegion   `json:"region"`
}

type sarifArtifact struct {
	URI string `json:"uri"`
}

type sarifRegion struct {
	StartLine   int `json:"startLine"`
	StartColumn int `json:"startColumn"`
}

// ReportSARIF writes every finding across all files as one SARIF 2.1.0
// document: a single run whose results each carry a physical location
// (file URI + 1-based line/column) so GitHub code scanning can ingest it.
func ReportSARIF(w io.Writer, files []FileViolations, toolVersion string) error {
	results := []sarifResult{}
	for _, f := range files {
		for _, v := range f.Violations {
			results = append(results, sarifResult{
				RuleID:  v.KataID,
				Level:   sarifLevel(v.Level),
				Message: sarifMessage{Text: v.Message},
				Locations: []sarifLocation{{
					PhysicalLocation: sarifPhysical{
						ArtifactLocation: sarifArtifact{URI: f.Filename},
						Region: sarifRegion{
							StartLine:   atLeastOne(v.Line),
							StartColumn: atLeastOne(v.Column),
						},
					},
				}},
			})
		}
	}
	doc := sarifDoc{
		Schema:  "https://json.schemastore.org/sarif-2.1.0.json",
		Version: "2.1.0",
		Runs: []sarifRun{{
			Tool: sarifTool{Driver: sarifDriver{
				Name:           "zshellcheck",
				InformationURI: "https://github.com/afadesigns/zshellcheck",
				Version:        toolVersion,
			}},
			Results: results,
		}},
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(doc)
}

// sarifLevel maps a kata severity onto a SARIF result level. SARIF has no
// "style"/"info" levels, so both map to "note".
func sarifLevel(s katas.Severity) string {
	switch s {
	case katas.SeverityError:
		return "error"
	case katas.SeverityWarning:
		return "warning"
	default:
		return "note"
	}
}

// atLeastOne clamps a 1-based SARIF coordinate so a zero never produces an
// invalid region.
func atLeastOne(n int) int {
	if n < 1 {
		return 1
	}
	return n
}
