// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package reporter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/config"
	"github.com/afadesigns/zshellcheck/pkg/katas"
)

func TestTextReporter_StyleSeverity(t *testing.T) {
	violations := []katas.Violation{
		{KataID: "ZC0001", Message: "style msg", Level: katas.SeverityStyle, Line: 1, Column: 1},
	}

	var buf bytes.Buffer
	cfg := config.DefaultConfig()
	reporter := NewTextReporter(&buf, "test.zsh", "echo hello", cfg)
	if err := reporter.Report(violations); err != nil {
		t.Fatalf("Report() error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "style msg") {
		t.Error("expected style message in output")
	}
	// Style severity uses cyan color
	if !strings.Contains(output, "\033[36m") {
		t.Error("expected cyan color code for style severity")
	}
}
