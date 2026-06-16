// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package reporter

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
)

func twoFiles() []FileViolations {
	return []FileViolations{
		{Filename: "a.zsh", Violations: []katas.Violation{
			{KataID: "ZC1001", Message: "msg a", Line: 1, Column: 1, Level: katas.SeverityError},
		}},
		{Filename: "b.zsh", Violations: []katas.Violation{
			{KataID: "ZC1002", Message: "msg b", Line: 5, Column: 10, Level: katas.SeverityWarning},
			{KataID: "ZC1003", Message: "msg c", Line: 7, Column: 2, Level: katas.SeverityStyle},
		}},
	}
}

func TestReportJSON_MultiFileIsOneValidArray(t *testing.T) {
	var buf bytes.Buffer
	if err := ReportJSON(&buf, twoFiles()); err != nil {
		t.Fatalf("ReportJSON error: %v", err)
	}
	// The whole output must be a single valid JSON array — the bug was one
	// array per file, concatenated.
	var findings []map[string]any
	if err := json.Unmarshal(buf.Bytes(), &findings); err != nil {
		t.Fatalf("multi-file JSON is not a single valid array: %v\n%s", err, buf.String())
	}
	if len(findings) != 3 {
		t.Fatalf("want 3 findings, got %d", len(findings))
	}
	if findings[0]["File"] != "a.zsh" || findings[1]["File"] != "b.zsh" {
		t.Errorf("findings not attributed to files: %v", findings)
	}
	if findings[0]["KataID"] != "ZC1001" {
		t.Errorf("want KataID ZC1001, got %v", findings[0]["KataID"])
	}
}

func TestReportJSON_EmptyIsValidArray(t *testing.T) {
	var buf bytes.Buffer
	if err := ReportJSON(&buf, nil); err != nil {
		t.Fatalf("ReportJSON error: %v", err)
	}
	var findings []any
	if err := json.Unmarshal(buf.Bytes(), &findings); err != nil {
		t.Fatalf("empty output is not valid JSON: %v", err)
	}
	if len(findings) != 0 {
		t.Errorf("want empty array, got %d", len(findings))
	}
}

func TestReportJSON_WriterError(t *testing.T) {
	if err := ReportJSON(&failWriter{}, twoFiles()); err == nil {
		t.Error("expected error from failing writer")
	}
}

func TestReportSARIF_MultiFileIsOneValidRun(t *testing.T) {
	var buf bytes.Buffer
	if err := ReportSARIF(&buf, twoFiles(), "1.2.3", testMeta); err != nil {
		t.Fatalf("ReportSARIF error: %v", err)
	}
	var doc map[string]any
	if err := json.Unmarshal(buf.Bytes(), &doc); err != nil {
		t.Fatalf("multi-file SARIF is not a single valid document: %v\n%s", err, buf.String())
	}
	if doc["version"] != "2.1.0" {
		t.Errorf("want SARIF 2.1.0, got %v", doc["version"])
	}
	runs := doc["runs"].([]any)
	if len(runs) != 1 {
		t.Fatalf("want exactly 1 run, got %d", len(runs))
	}
	run := runs[0].(map[string]any)
	driver := run["tool"].(map[string]any)["driver"].(map[string]any)
	if driver["name"] != "zshellcheck" || driver["version"] != "1.2.3" {
		t.Errorf("driver metadata wrong: %v", driver)
	}
	results := run["results"].([]any)
	if len(results) != 3 {
		t.Fatalf("want 3 results, got %d", len(results))
	}
	// First result carries a physical location with the file and 1-based
	// line/column.
	r0 := results[0].(map[string]any)
	if r0["level"] != "error" {
		t.Errorf("want level error, got %v", r0["level"])
	}
	loc := r0["locations"].([]any)[0].(map[string]any)["physicalLocation"].(map[string]any)
	if loc["artifactLocation"].(map[string]any)["uri"] != "a.zsh" {
		t.Errorf("want uri a.zsh, got %v", loc["artifactLocation"])
	}
	region := loc["region"].(map[string]any)
	if region["startLine"].(float64) != 1 || region["startColumn"].(float64) != 1 {
		t.Errorf("want region 1:1, got %v", region)
	}
}

func TestReportSARIF_LevelMappingAndClamp(t *testing.T) {
	files := []FileViolations{{Filename: "z.zsh", Violations: []katas.Violation{
		{KataID: "ZC1", Message: "i", Line: 0, Column: 0, Level: katas.SeverityInfo},
		{KataID: "ZC2", Message: "s", Line: 1, Column: 1, Level: katas.SeverityStyle},
	}}}
	var buf bytes.Buffer
	// nil meta exercises the minimal-descriptor path.
	if err := ReportSARIF(&buf, files, "0.0.0", nil); err != nil {
		t.Fatalf("ReportSARIF error: %v", err)
	}
	var doc map[string]any
	if err := json.Unmarshal(buf.Bytes(), &doc); err != nil {
		t.Fatalf("not valid SARIF: %v", err)
	}
	results := doc["runs"].([]any)[0].(map[string]any)["results"].([]any)
	// Info and Style both map to "note".
	for _, r := range results {
		if r.(map[string]any)["level"] != "note" {
			t.Errorf("want note level, got %v", r.(map[string]any)["level"])
		}
	}
	// A zero line/column is clamped to 1 so the SARIF region stays valid.
	region := results[0].(map[string]any)["locations"].([]any)[0].(map[string]any)["physicalLocation"].(map[string]any)["region"].(map[string]any)
	if region["startLine"].(float64) != 1 || region["startColumn"].(float64) != 1 {
		t.Errorf("zero coords not clamped to 1: %v", region)
	}
}

func TestReportSARIF_RulesMetadata(t *testing.T) {
	var buf bytes.Buffer
	if err := ReportSARIF(&buf, twoFiles(), "1.2.3", testMeta); err != nil {
		t.Fatalf("ReportSARIF error: %v", err)
	}
	var doc map[string]any
	if err := json.Unmarshal(buf.Bytes(), &doc); err != nil {
		t.Fatalf("not valid SARIF: %v", err)
	}
	run := doc["runs"].([]any)[0].(map[string]any)
	// One rule per distinct kata, with description, help URI, and level.
	rules := run["tool"].(map[string]any)["driver"].(map[string]any)["rules"].([]any)
	if len(rules) != 3 {
		t.Fatalf("want 3 rules, got %d", len(rules))
	}
	rule0 := rules[0].(map[string]any)
	if rule0["id"] != "ZC1001" {
		t.Errorf("rule0 id = %v", rule0["id"])
	}
	if rule0["helpUri"] == "" {
		t.Error("rule0 missing helpUri")
	}
	if rule0["fullDescription"].(map[string]any)["text"] != "ZC1001 desc" {
		t.Errorf("rule0 description = %v", rule0["fullDescription"])
	}
	// Results reference their rule by index.
	results := run["results"].([]any)
	if results[1].(map[string]any)["ruleIndex"].(float64) != 1 {
		t.Errorf("ruleIndex not wired: %v", results[1])
	}
}

func testMeta(id string) RuleMeta {
	return RuleMeta{
		Name:        id + "-name",
		Title:       id + " title",
		Description: id + " desc",
		HelpURI:     "https://example.test/" + id,
	}
}

func TestReportSARIF_WriterError(t *testing.T) {
	if err := ReportSARIF(&failWriter{}, twoFiles(), "1.0.0", testMeta); err == nil {
		t.Error("expected error from failing writer")
	}
}
