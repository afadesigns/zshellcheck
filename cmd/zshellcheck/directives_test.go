// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/config"
	"github.com/afadesigns/zshellcheck/pkg/katas"
)

func TestParseRuleSeverity(t *testing.T) {
	m, err := parseRuleSeverity("ZC1037:error, ZC1075:style ,")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["ZC1037"] != katas.SeverityError || m["ZC1075"] != katas.SeverityStyle {
		t.Errorf("re-grade map wrong: %v", m)
	}
	if got, err := parseRuleSeverity(""); got != nil || err != nil {
		t.Errorf("empty spec should be nil/nil, got %v/%v", got, err)
	}
	for _, bad := range []string{"ZC1037", "ZC1037:bogus", ":error"} {
		if _, err := parseRuleSeverity(bad); err == nil {
			t.Errorf("expected error for %q", bad)
		}
	}
}

func TestRegradeSeverity(t *testing.T) {
	vs := []katas.Violation{
		{KataID: "ZC1037", Level: katas.SeverityStyle},
		{KataID: "ZC1075", Level: katas.SeverityWarning},
	}
	regradeSeverity(vs, map[string]katas.Severity{"ZC1037": katas.SeverityError})
	if vs[0].Level != katas.SeverityError {
		t.Errorf("ZC1037 not re-graded: %v", vs[0].Level)
	}
	if vs[1].Level != katas.SeverityWarning {
		t.Errorf("ZC1075 should be unchanged: %v", vs[1].Level)
	}
	// Empty map is a no-op.
	regradeSeverity(vs, nil)
	if vs[0].Level != katas.SeverityError {
		t.Error("nil map should not change levels")
	}
}

func TestReportStaleNoka(t *testing.T) {
	raw := []katas.Violation{{KataID: "ZC1002", Line: 1}}
	directives := config.Directives{PerLine: map[int][]string{
		1: {"ZC1002", "ZC9999", "ZC9998"}, // ZC1002 fires; two stale on one line
		2: {"ZC1037"},                     // stale: nothing on line 2
	}}
	var buf bytes.Buffer
	count := 0
	reportStaleNoka(&buf, "f.zsh", raw, directives, &count)
	out := buf.String()
	if count != 3 {
		t.Errorf("want 3 stale, got %d:\n%s", count, out)
	}
	if strings.Contains(out, "ZC1002") {
		t.Error("ZC1002 fires and must not be reported stale")
	}
	// Sorted by line, then ID: ZC9998 < ZC9999 (same line 1) < ZC1037 (line 2).
	if !(strings.Index(out, "ZC9998") < strings.Index(out, "ZC9999") &&
		strings.Index(out, "ZC9999") < strings.Index(out, "ZC1037")) {
		t.Errorf("stale output not sorted by line then id:\n%s", out)
	}
}

func TestAddNokaDirectives(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "t.zsh")
	src := "echo a\nx=1  # noka: ZC1\necho b\n"
	if err := os.WriteFile(path, []byte(src), 0o600); err != nil {
		t.Fatal(err)
	}
	vs := []katas.Violation{
		{KataID: "ZC1002", Line: 1},
		{KataID: "ZC1037", Line: 1}, // two on line 1 -> grouped
		{KataID: "ZC9", Line: 2},    // line already has a noka -> skipped
	}
	if err := addNokaDirectives(path, []byte(src), vs); err != nil {
		t.Fatalf("addNokaDirectives: %v", err)
	}
	out, _ := os.ReadFile(path)
	lines := strings.Split(string(out), "\n")
	if lines[0] != "echo a  # noka: ZC1002, ZC1037" {
		t.Errorf("line 1 = %q", lines[0])
	}
	if lines[1] != "x=1  # noka: ZC1" {
		t.Errorf("line 2 (existing noka) should be untouched: %q", lines[1])
	}

	// No findings -> no write, no error.
	if err := addNokaDirectives(path, []byte(src), nil); err != nil {
		t.Errorf("empty add-noka should be a no-op: %v", err)
	}
}

func TestReportStaleNokaNilCounter(t *testing.T) {
	d := config.Directives{PerLine: map[int][]string{1: {"ZC1"}}}
	var buf bytes.Buffer
	reportStaleNoka(&buf, "f.zsh", nil, d, nil) // nil counter must not panic
	if !strings.Contains(buf.String(), "ZC1") {
		t.Error("stale entry should still be reported with a nil counter")
	}
}

func TestAddNokaDirectivesOutOfRange(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "t.zsh")
	src := "echo a\n"
	if err := os.WriteFile(path, []byte(src), 0o600); err != nil {
		t.Fatal(err)
	}
	// A finding on a line past EOF is skipped; the file is left unchanged.
	if err := addNokaDirectives(path, []byte(src), []katas.Violation{{KataID: "ZC1", Line: 99}}); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(path); string(b) != src {
		t.Errorf("out-of-range finding modified the file: %q", b)
	}
}

func TestRunBaselineWriteError(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "t.zsh")
	if err := os.WriteFile(src, []byte("echo hi\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	old := os.Args
	defer func() { os.Args = old }()
	resetFlags()
	// A baseline path inside a non-existent directory cannot be written.
	os.Args = []string{"zshellcheck", "-no-banner", "-baseline-write", filepath.Join(dir, "nope", "b.txt"), src}
	if code := run(); code != 1 {
		t.Errorf("baseline write to a bad path should exit 1, got %d", code)
	}
}

func TestConfigureModes(t *testing.T) {
	resetFlags()
	flags := registerRunFlags()
	*flags.statistics = true
	*flags.detectStale = true
	*flags.addNoka = true
	*flags.ruleSeverity = "ZC1037:error"
	var opts fixOptions
	if code := configureModes(flags, &opts); code != 0 {
		t.Fatalf("configureModes code=%d", code)
	}
	if opts.statistics == nil || !opts.detectStale || !opts.addNoka || opts.staleCount == nil {
		t.Errorf("modes not configured: %+v", opts)
	}
	if opts.ruleSeverity["ZC1037"] != katas.SeverityError {
		t.Error("rule-severity not parsed")
	}
	// Malformed rule-severity -> exit 1.
	resetFlags()
	flags2 := registerRunFlags()
	*flags2.ruleSeverity = "bogus"
	var opts2 fixOptions
	if code := configureModes(flags2, &opts2); code != 1 {
		t.Errorf("malformed rule-severity should exit 1, got %d", code)
	}
}

func TestRun_RuleSeverityAndStale(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "t.zsh")
	old := os.Args
	defer func() { os.Args = old }()

	// rule-severity re-grade still exits 1 on findings.
	if err := os.WriteFile(path, []byte("echo hi\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	resetFlags()
	os.Args = []string{"zshellcheck", "-no-banner", "-rule-severity", "ZC1037:error", path}
	if code := run(); code != 1 {
		t.Errorf("findings should exit 1, got %d", code)
	}

	// A stale directive on an otherwise clean file exits 1 under detection.
	if err := os.WriteFile(path, []byte(": # noka: ZC9999\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	resetFlags()
	os.Args = []string{"zshellcheck", "-no-banner", "-detect-stale-noka", path}
	if code := run(); code != 1 {
		t.Errorf("stale directive should exit 1, got %d", code)
	}

	// add-noka exits 0 and writes the directives.
	if err := os.WriteFile(path, []byte("echo hi\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	resetFlags()
	os.Args = []string{"zshellcheck", "-no-banner", "-add-noka", path}
	if code := run(); code != 0 {
		t.Errorf("add-noka should exit 0, got %d", code)
	}
	if b, _ := os.ReadFile(path); !strings.Contains(string(b), "# noka:") {
		t.Errorf("add-noka did not write directive: %q", b)
	}
}
