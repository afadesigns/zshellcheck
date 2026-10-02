// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package config

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
)

func TestParseBlockList(t *testing.T) {
	cfg, err := Parse([]byte("no_color: true\ndisabled_katas:\n  - ZC1001\n  - ZC1002\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.NoColor {
		t.Error("want NoColor true")
	}
	if cfg.ModifiedSeverities["ZC1001"] != katas.SeverityDisabled {
		t.Errorf("ZC1001 = %q, want disabled", cfg.ModifiedSeverities["ZC1001"])
	}
	if cfg.ModifiedSeverities["ZC1002"] != katas.SeverityDisabled {
		t.Errorf("ZC1002 = %q, want disabled", cfg.ModifiedSeverities["ZC1002"])
	}
}

func TestParseInlineList(t *testing.T) {
	cfg, err := Parse([]byte("disabled_katas: [ZC1, 'ZC2', \"ZC3\"]\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, id := range []string{"ZC1", "ZC2", "ZC3"} {
		if cfg.ModifiedSeverities[id] != katas.SeverityDisabled {
			t.Errorf("%s = %q, want disabled", id, cfg.ModifiedSeverities[id])
		}
	}
}

func TestParseEmptyInlineList(t *testing.T) {
	cfg, err := Parse([]byte("disabled_katas: []\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.ModifiedSeverities) != 0 {
		t.Errorf("want empty, got %v", cfg.ModifiedSeverities)
	}
}

func TestParseSeverityBlockList(t *testing.T) {
	cfg, err := Parse([]byte("kata_severity:\n  error:\n    - ZC1001\n  warning:\n    - ZC1408\n  disabled:\n    - ZC9999\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ModifiedSeverities["ZC1001"] != katas.SeverityError {
		t.Errorf("ZC1001 = %q, want error", cfg.ModifiedSeverities["ZC1001"])
	}
	if cfg.ModifiedSeverities["ZC1408"] != katas.SeverityWarning {
		t.Errorf("ZC1408 = %q, want warning", cfg.ModifiedSeverities["ZC1408"])
	}
	if cfg.ModifiedSeverities["ZC9999"] != katas.SeverityDisabled {
		t.Errorf("ZC9999 = %q, want disabled", cfg.ModifiedSeverities["ZC9999"])
	}
}

func TestParseDisabledKatasBackwardCompat(t *testing.T) {
	cfg, err := Parse([]byte("disabled_katas:\n  - ZC1001\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ModifiedSeverities["ZC1001"] != katas.SeverityDisabled {
		t.Errorf("ZC1001 = %q, want disabled", cfg.ModifiedSeverities["ZC1001"])
	}
}

func TestParseAllScalars(t *testing.T) {
	src := "error_color: a\nwarning_color: b\ninfo_color: c\nid_color: d\n" +
		"title_color: e\nmessage_color: f\nline_color: g\ncolumn_color: h\n" +
		"no_color: false\nverbose: true\n"
	cfg, err := Parse([]byte(src))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := []string{
		cfg.ErrorColor, cfg.WarningColor, cfg.InfoColor, cfg.IDColor,
		cfg.TitleColor, cfg.MessageColor, cfg.LineColor, cfg.ColumnColor,
	}
	want := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("color[%d] = %q, want %q", i, got[i], want[i])
		}
	}
	if cfg.NoColor {
		t.Error("want NoColor false")
	}
	if !cfg.Verbose {
		t.Error("want Verbose true")
	}
}

func TestParseQuotesAndEscapes(t *testing.T) {
	cfg, err := Parse([]byte("error_color: \"\\e[31m\"\nwarning_color: '\\eliteral'\n" +
		"info_color: \"\\x1b[1m\"\nid_color: \"\\x1b[0m\"\ntitle_color: \"tab\\there\"\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ErrorColor != "\x1b[31m" {
		t.Errorf("ErrorColor = %q, want ESC[31m", cfg.ErrorColor)
	}
	if cfg.WarningColor != "\\eliteral" {
		t.Errorf("single-quote should be literal, got %q", cfg.WarningColor)
	}
	if cfg.InfoColor != "\x1b[1m" {
		t.Errorf("InfoColor = %q, want ESC[1m", cfg.InfoColor)
	}
	if cfg.IDColor != "\x1b[0m" {
		t.Errorf("IDColor = %q, want ESC[0m", cfg.IDColor)
	}
	if cfg.TitleColor != "tab\there" {
		t.Errorf("TitleColor = %q, want tab<TAB>here", cfg.TitleColor)
	}
}

func TestParseComments(t *testing.T) {
	cfg, err := Parse([]byte("# full line\nno_color: true  # trailing\nerror_color: \"#ff0000\"\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.NoColor {
		t.Error("want NoColor true")
	}
	if cfg.ErrorColor != "#ff0000" {
		t.Errorf("quoted hash should survive, got %q", cfg.ErrorColor)
	}
}

func TestParseUnknownKeyIgnored(t *testing.T) {
	cfg, err := Parse([]byte("future_option: whatever\nno_color: true\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.NoColor {
		t.Error("want NoColor true after ignored unknown key")
	}
}

func TestParseErrors(t *testing.T) {
	for _, src := range []string{
		":this is not yaml\n",
		"  - orphan\n",
		"no_color: maybe\n",
		"verbose: notabool\n",
		"bareword\n",
	} {
		if _, err := Parse([]byte(src)); err == nil {
			t.Errorf("expected error for %q", src)
		}
	}
}

func TestParseEmptyAndBlankLines(t *testing.T) {
	cfg, err := Parse([]byte("\n\n   \n# only comments\n\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.ModifiedSeverities) != 0 {
		t.Errorf("blank config should have empty ModifiedSeverities, got %+v", cfg.ModifiedSeverities)
	}
}

func TestParseKataSeverityFlowStyle(t *testing.T) {
	cfg, warnings, err := ParseWithWarnings([]byte("kata_severity: {error: [ZC1], style: [ZC2]}\n"))
	if err != nil || len(warnings) != 0 {
		t.Fatalf("err=%v warnings=%v", err, warnings)
	}
	if cfg.ModifiedSeverities["ZC1"] != katas.SeverityError || cfg.ModifiedSeverities["ZC2"] != katas.SeverityStyle {
		t.Errorf("got %v", cfg.ModifiedSeverities)
	}
}

func TestParseUnknownKataSeverityKeyWarns(t *testing.T) {
	cfg, warnings, err := ParseWithWarnings([]byte("kata_severity:\n  eror:\n    - ZC1\n  error:\n    - ZC2\n"))
	if err != nil || len(warnings) != 1 {
		t.Fatalf("err=%v warnings=%v, want one warning", err, warnings)
	}
	if _, ok := cfg.ModifiedSeverities["ZC1"]; ok || cfg.ModifiedSeverities["ZC2"] != katas.SeverityError {
		t.Errorf("got %v", cfg.ModifiedSeverities)
	}
}

func TestParseRepeatedIDHighestWinsWithWarning(t *testing.T) {
	src := "disabled_katas: [ZC1, ZC3]\nkata_severity:\n  error: [ZC1]\n  style: [ZC2, ZC2]\n  info: [ZC3]\n"
	cfg, warnings, err := ParseWithWarnings([]byte(src))
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]katas.Severity{"ZC1": katas.SeverityError, "ZC2": katas.SeverityStyle, "ZC3": katas.SeverityInfo}
	for id, sev := range want {
		if cfg.ModifiedSeverities[id] != sev {
			t.Errorf("%s = %q, want %q", id, cfg.ModifiedSeverities[id], sev)
		}
	}
	if len(warnings) != 3 {
		t.Errorf("want 3 warnings, got %v", warnings)
	}
}

func TestParseWrongShapeErrors(t *testing.T) {
	for _, src := range []string{
		"kata_severity: foo\n",
		"kata_severity:\n  error: ZC1\n",
		"disabled_katas: ZC1\n",
		"kata_severity:\n  error:\n    - [a]\n",
	} {
		if _, err := Parse([]byte(src)); err == nil {
			t.Errorf("expected error for %q", src)
		}
	}
}

func TestParseCRLF(t *testing.T) {
	cfg, err := Parse([]byte("no_color: true\r\nkata_severity:\r\n  error:\r\n    - ZC1\r\n"))
	if err != nil || !cfg.NoColor || cfg.ModifiedSeverities["ZC1"] != katas.SeverityError {
		t.Errorf("err=%v cfg=%+v", err, cfg)
	}
}

func TestParseTabIndentIsRejected(t *testing.T) {
	if _, err := Parse([]byte("kata_severity:\n\terror:\n\t\t- ZC1\n")); err == nil {
		t.Error("YAML forbids tab indentation; want error")
	}
}
