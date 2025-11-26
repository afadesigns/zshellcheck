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
	return &SarifReporter{writer: writer, filename: filename}
}

// Report prints the violations to the writer in SARIF format.
func (r *SarifReporter) Report(violations []katas.Violation) error {
	run := SarifRun{
		Tool: SarifTool{
			Driver: SarifDriver{
				Name:            "ZShellCheck",
				InformationUri:  "https://github.com/afadesigns/zshellcheck",
				SemanticVersion: "0.0.92", // TODO: Use version variable
				Rules:           make([]SarifRule, 0),
			},
		},
		Results: make([]SarifResult, 0),
	}

	// Map to track unique rules
	rulesMap := make(map[string]SarifRule)

	for _, v := range violations {
		kata, ok := katas.Registry.GetKata(v.KataID)
		if !ok {
			continue
		}

		if _, exists := rulesMap[v.KataID]; !exists {
			rulesMap[v.KataID] = SarifRule{
				Id:               v.KataID,
				Name:             v.KataID,
				ShortDescription: SarifMessage{Text: kata.Title},
				FullDescription:  SarifMessage{Text: kata.Description},
				HelpUri:          "https://github.com/afadesigns/zshellcheck/blob/main/KATAS.md#" + v.KataID,
			}
		}

		result := SarifResult{
			RuleId:  v.KataID,
			Level:   "warning", // Default level
			Message: SarifMessage{Text: v.Message},
			Locations: []SarifLocation{
				{
					PhysicalLocation: SarifPhysicalLocation{
						ArtifactLocation: SarifArtifactLocation{
							Uri: r.filename,
						},
						Region: SarifRegion{
							StartLine:   v.Line,
							StartColumn: v.Column,
						},
					},
				},
			},
		}
		run.Results = append(run.Results, result)
	}

	for _, rule := range rulesMap {
		run.Tool.Driver.Rules = append(run.Tool.Driver.Rules, rule)
	}

	report := SarifReport{
		Version: "2.1.0",
		Schema:  "https://schemastore.azurewebsites.net/schemas/json/sarif-2.1.0-rtm.5.json",
		Runs:    []SarifRun{run},
	}

	enc := json.NewEncoder(r.writer)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}

type SarifReport struct {
	Version string     `json:"version"`
	Schema  string     `json:"$schema"`
	Runs    []SarifRun `json:"runs"`
}

type SarifRun struct {
	Tool    SarifTool     `json:"tool"`
	Results []SarifResult `json:"results"`
}

type SarifTool struct {
	Driver SarifDriver `json:"driver"`
}

type SarifDriver struct {
	Name            string      `json:"name"`
	InformationUri  string      `json:"informationUri"`
	SemanticVersion string      `json:"semanticVersion"`
	Rules           []SarifRule `json:"rules"`
}

type SarifRule struct {
	Id               string       `json:"id"`
	Name             string       `json:"name"`
	ShortDescription SarifMessage `json:"shortDescription"`
	FullDescription  SarifMessage `json:"fullDescription"`
	HelpUri          string       `json:"helpUri"`
}

type SarifResult struct {
	RuleId    string          `json:"ruleId"`
	Level     string          `json:"level"` // error, warning, note
	Message   SarifMessage    `json:"message"`
	Locations []SarifLocation `json:"locations"`
}

type SarifMessage struct {
	Text string `json:"text"`
}

type SarifLocation struct {
	PhysicalLocation SarifPhysicalLocation `json:"physicalLocation"`
}

type SarifPhysicalLocation struct {
	ArtifactLocation SarifArtifactLocation `json:"artifactLocation"`
	Region           SarifRegion           `json:"region"`
}

type SarifArtifactLocation struct {
	Uri string `json:"uri"`
}

type SarifRegion struct {
	StartLine   int `json:"startLine"`
	StartColumn int `json:"startColumn"`
}
