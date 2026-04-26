// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package main

import (
	"reflect"
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/config"
	"github.com/afadesigns/zshellcheck/pkg/katas"
)

func TestMergeDisabledNoExtra(t *testing.T) {
	base := []string{"ZC1001", "ZC1002"}
	got := mergeDisabled(base, nil)
	if !reflect.DeepEqual(got, base) {
		t.Errorf("expected base, got %v", got)
	}
}

func TestMergeDisabledExtraAppends(t *testing.T) {
	base := []string{"ZC1001"}
	got := mergeDisabled(base, []string{"ZC1002", "ZC1003"})
	want := []string{"ZC1001", "ZC1002", "ZC1003"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestApplyDirectiveSilencesNoDirectives(t *testing.T) {
	violations := []katas.Violation{{KataID: "ZC1001", Line: 1}}
	edits := []katas.FixEdit{{Line: 1, Column: 1, Length: 1, Replace: "x"}}
	gotV, gotE := applyDirectiveSilences(violations, edits, config.Directives{})
	if len(gotV) != 1 || len(gotE) != 1 {
		t.Errorf("expected pass-through, got vs=%d es=%d", len(gotV), len(gotE))
	}
}

func TestApplyDirectiveSilencesPerLineDrop(t *testing.T) {
	violations := []katas.Violation{
		{KataID: "ZC1001", Line: 5},
		{KataID: "ZC1002", Line: 6},
	}
	edits := []katas.FixEdit{}
	directives := config.Directives{
		PerLine: map[int][]string{5: {"ZC1001"}},
	}
	gotV, _ := applyDirectiveSilences(violations, edits, directives)
	if len(gotV) != 1 || gotV[0].KataID != "ZC1002" {
		t.Errorf("expected ZC1002 only, got %v", gotV)
	}
}

func TestApplyDirectiveSilencesPerLineAll(t *testing.T) {
	violations := []katas.Violation{
		{KataID: "ZC1001", Line: 5},
		{KataID: "ZC1002", Line: 5},
	}
	directives := config.Directives{
		PerLineAll: map[int]bool{5: true},
		// applyDirectiveSilences gates on PerLine being non-empty —
		// give it one entry so the silence path runs and the
		// PerLineAll lookup inside IsDisabledOn fires for both
		// violations on the bare-noka line.
		PerLine: map[int][]string{5: {"ZC1001"}},
	}
	gotV, _ := applyDirectiveSilences(violations, nil, directives)
	if len(gotV) != 0 {
		t.Errorf("expected all silenced, got %v", gotV)
	}
}

func TestApplySeverityFilterEmpty(t *testing.T) {
	violations := []katas.Violation{{KataID: "ZC1001", Level: katas.SeverityError}}
	gotV, _ := applySeverityFilter(violations, nil, nil)
	if len(gotV) != 1 {
		t.Errorf("expected pass-through, got %v", gotV)
	}
}

func TestApplySeverityFilterMatch(t *testing.T) {
	violations := []katas.Violation{
		{KataID: "ZC1001", Level: katas.SeverityError},
		{KataID: "ZC1002", Level: katas.SeverityStyle},
	}
	gotV, _ := applySeverityFilter(violations, nil, []katas.Severity{katas.SeverityError})
	if len(gotV) != 1 || gotV[0].KataID != "ZC1001" {
		t.Errorf("expected only ZC1001, got %v", gotV)
	}
}

func TestParseSeverityFilterValid(t *testing.T) {
	got, code := parseSeverityFilter("error,warning")
	if code != 0 {
		t.Fatalf("unexpected code %d", code)
	}
	if len(got) != 2 || got[0] != katas.SeverityError || got[1] != katas.SeverityWarning {
		t.Errorf("unexpected: %v", got)
	}
}

func TestParseSeverityFilterInvalid(t *testing.T) {
	_, code := parseSeverityFilter("frobozz")
	if code != 1 {
		t.Errorf("expected exit 1, got %d", code)
	}
}

func TestParseSeverityFilterEmpty(t *testing.T) {
	got, code := parseSeverityFilter("")
	if code != 0 || got != nil {
		t.Errorf("expected nil, got %v code %d", got, code)
	}
}

func TestBuildFixOptsAllOff(t *testing.T) {
	got := buildFixOpts(false, false, false)
	if got.enabled || got.diff || got.dryRun {
		t.Errorf("expected all-off, got %+v", got)
	}
}

func TestBuildFixOptsFix(t *testing.T) {
	got := buildFixOpts(true, false, false)
	if !got.enabled || got.diff || got.dryRun {
		t.Errorf("expected enabled-only, got %+v", got)
	}
	if got.stats == nil {
		t.Errorf("expected stats allocated")
	}
}

func TestBuildFixOptsDiff(t *testing.T) {
	got := buildFixOpts(false, true, false)
	if !got.enabled || !got.diff || !got.dryRun {
		t.Errorf("expected diff implies dry-run + enabled, got %+v", got)
	}
}

func TestBuildFixOptsDryRun(t *testing.T) {
	got := buildFixOpts(true, false, true)
	if !got.enabled || got.diff || !got.dryRun {
		t.Errorf("expected fix+dry-run, got %+v", got)
	}
}

func TestApplyFlagOverridesNoColor(t *testing.T) {
	cfg := config.DefaultConfig()
	got := applyFlagOverrides(cfg, true, false)
	if !got.NoColor {
		t.Error("expected NoColor=true")
	}
}

func TestApplyFlagOverridesVerbose(t *testing.T) {
	cfg := config.DefaultConfig()
	got := applyFlagOverrides(cfg, false, true)
	if !got.Verbose {
		t.Error("expected Verbose=true")
	}
}
