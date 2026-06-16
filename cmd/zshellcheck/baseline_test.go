// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
)

func TestBaselineFingerprint(t *testing.T) {
	lines := []string{"echo one", "echo two"}
	v := katas.Violation{KataID: "ZC1037", Line: 2, Column: 1}
	if got := baselineFingerprint("a.zsh", lines, v); got != "ZC1037\ta.zsh\techo two" {
		t.Errorf("fingerprint = %q", got)
	}
	// Out-of-range line yields an empty content field, never a panic.
	vOut := katas.Violation{KataID: "ZC1", Line: 99}
	if got := baselineFingerprint("a.zsh", lines, vOut); got != "ZC1\ta.zsh\t" {
		t.Errorf("out-of-range fingerprint = %q", got)
	}
}

func TestApplyBaselineWriteMode(t *testing.T) {
	b := &baselineState{write: true}
	vs := []katas.Violation{
		{KataID: "ZC1", Line: 1}, {KataID: "ZC2", Line: 1},
	}
	got := b.applyBaseline("f.zsh", []byte("cmd\n"), vs)
	if len(got) != 2 {
		t.Errorf("write mode should return all findings, got %d", len(got))
	}
	if len(b.collect) != 2 {
		t.Errorf("write mode should collect 2 fingerprints, got %d", len(b.collect))
	}
}

func TestApplyBaselineFilterMode(t *testing.T) {
	data := []byte("cmd\n")
	known := baselineFingerprint("f.zsh", []string{"cmd"}, katas.Violation{KataID: "ZC1", Line: 1})
	b := &baselineState{known: map[string]bool{known: true}}
	vs := []katas.Violation{
		{KataID: "ZC1", Line: 1}, // in baseline -> suppressed
		{KataID: "ZC2", Line: 1}, // new -> kept
	}
	got := b.applyBaseline("f.zsh", data, vs)
	if len(got) != 1 || got[0].KataID != "ZC2" {
		t.Errorf("filter mode kept wrong findings: %v", got)
	}
}

func TestWriteAndLoadBaseline(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "base.txt")
	b := &baselineState{write: true, collect: []string{
		"ZC2\tf\tline", "ZC1\tf\tline", "ZC2\tf\tline", // duplicate + unsorted
	}}
	if err := b.writeBaseline(path); err != nil {
		t.Fatalf("writeBaseline: %v", err)
	}
	data, _ := os.ReadFile(path)
	got := strings.TrimSpace(string(data))
	// De-duplicated and sorted.
	if got != "ZC1\tf\tline\nZC2\tf\tline" {
		t.Errorf("baseline content = %q", got)
	}
	loaded, err := loadBaseline(path)
	if err != nil {
		t.Fatalf("loadBaseline: %v", err)
	}
	if !loaded.known["ZC1\tf\tline"] || !loaded.known["ZC2\tf\tline"] {
		t.Errorf("loaded baseline missing entries: %v", loaded.known)
	}
}

func TestLoadBaselineMissing(t *testing.T) {
	if _, err := loadBaseline("/nonexistent/baseline.txt"); err == nil {
		t.Error("expected error for missing baseline file")
	}
}

func TestSetupBaseline(t *testing.T) {
	// Write mode.
	var opts fixOptions
	if code := setupBaseline(&opts, "", "out.txt"); code != 0 || opts.baseline == nil || !opts.baseline.write {
		t.Errorf("write-mode setup wrong: code=%d baseline=%v", code, opts.baseline)
	}
	// Filter mode with a real file.
	dir := t.TempDir()
	path := filepath.Join(dir, "b.txt")
	if err := os.WriteFile(path, []byte("ZC1\tf\tline\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	opts = fixOptions{}
	if code := setupBaseline(&opts, path, ""); code != 0 || opts.baseline == nil || opts.baseline.write {
		t.Errorf("filter-mode setup wrong: code=%d", code)
	}
	// Missing file errors.
	opts = fixOptions{}
	if code := setupBaseline(&opts, "/nonexistent/zzz.txt", ""); code != 1 {
		t.Errorf("missing baseline should exit 1, got %d", code)
	}
	// Neither flag set leaves baseline nil.
	opts = fixOptions{}
	if code := setupBaseline(&opts, "", ""); code != 0 || opts.baseline != nil {
		t.Errorf("no-baseline setup should be inert, got code=%d", code)
	}
}

func TestRun_BaselineRatchet(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "t.zsh")
	if err := os.WriteFile(src, []byte("x=`which git`\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	base := filepath.Join(dir, "base.txt")
	old := os.Args
	defer func() { os.Args = old }()

	// Write a baseline of the current findings; exit 0.
	resetFlags()
	os.Args = []string{"zshellcheck", "-no-banner", "-baseline-write", base, src}
	if code := run(); code != 0 {
		t.Errorf("baseline-write should exit 0, got %d", code)
	}
	if _, err := os.Stat(base); err != nil {
		t.Fatalf("baseline file not written: %v", err)
	}

	// With the baseline, the same source yields no new findings: exit 0.
	resetFlags()
	os.Args = []string{"zshellcheck", "-no-banner", "-baseline", base, src}
	if code := run(); code != 0 {
		t.Errorf("all-baselined run should exit 0, got %d", code)
	}

	// A new finding fails the ratchet: exit 1.
	if err := os.WriteFile(src, []byte("x=`which git`\nrm -rf $undefined\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	resetFlags()
	os.Args = []string{"zshellcheck", "-no-banner", "-baseline", base, src}
	if code := run(); code != 1 {
		t.Errorf("new finding should exit 1, got %d", code)
	}
}
