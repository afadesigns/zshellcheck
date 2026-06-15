// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/katas"
)

func testRulesRegistry() *katas.KatasRegistry {
	kr := katas.NewKatasRegistry()
	node := &ast.SimpleCommand{}
	kr.RegisterKata(node, katas.Kata{
		ID: "ZC1002", Title: "Bravo title", Description: "Bravo body.", Severity: katas.SeverityWarning,
	})
	kr.RegisterKata(node, katas.Kata{
		ID: "ZC1001", Title: "Alpha title", Description: "Alpha body.", Severity: katas.SeverityStyle,
	})
	return kr
}

func TestPrintRulesList(t *testing.T) {
	var out bytes.Buffer
	if code := printRulesList(&out, testRulesRegistry()); code != 0 {
		t.Fatalf("printRulesList code = %d, want 0", code)
	}
	got := out.String()
	for _, want := range []string{"ZC1001", "ZC1002", "Alpha title", "Bravo title", "2 katas."} {
		if !strings.Contains(got, want) {
			t.Errorf("--list-rules output missing %q\n%s", want, got)
		}
	}
	// IDs are sorted: ZC1001 precedes ZC1002.
	if i, j := strings.Index(got, "ZC1001"), strings.Index(got, "ZC1002"); i > j {
		t.Errorf("--list-rules not sorted by ID: %s", got)
	}
}

func TestPrintRuleExplainFound(t *testing.T) {
	var out, errOut bytes.Buffer
	for _, id := range []string{"ZC1001", "zc1001", " ZC1001 "} {
		out.Reset()
		errOut.Reset()
		if code := printRuleExplain(&out, &errOut, testRulesRegistry(), id); code != 0 {
			t.Fatalf("explain %q code = %d, want 0", id, code)
		}
		got := out.String()
		for _, want := range []string{"ZC1001", "Alpha title", "Severity: Style", "Alpha body."} {
			if !strings.Contains(got, want) {
				t.Errorf("explain %q output missing %q\n%s", id, want, got)
			}
		}
	}
}

func TestPrintRuleExplainUnknown(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := printRuleExplain(&out, &errOut, testRulesRegistry(), "ZC9999"); code != 1 {
		t.Fatalf("explain unknown code = %d, want 1", code)
	}
	if out.Len() != 0 {
		t.Errorf("explain unknown wrote to stdout: %q", out.String())
	}
	if !strings.Contains(errOut.String(), "ZC9999") {
		t.Errorf("explain unknown stderr missing ID: %q", errOut.String())
	}
}

func TestTitleSeverity(t *testing.T) {
	cases := map[katas.Severity]string{
		katas.SeverityError: "Error", katas.SeverityWarning: "Warning",
		katas.SeverityInfo: "Info", katas.SeverityStyle: "Style", "": "",
	}
	for in, want := range cases {
		if got := titleSeverity(in); got != want {
			t.Errorf("titleSeverity(%q) = %q, want %q", in, got, want)
		}
	}
}
