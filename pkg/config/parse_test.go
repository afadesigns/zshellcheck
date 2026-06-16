// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package config

import (
	"reflect"
	"testing"
)

func TestParseBlockList(t *testing.T) {
	cfg, err := Parse([]byte("no_color: true\ndisabled_katas:\n  - ZC1001\n  - ZC1002\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.NoColor {
		t.Error("want NoColor true")
	}
	if want := []string{"ZC1001", "ZC1002"}; !reflect.DeepEqual(cfg.DisabledKatas, want) {
		t.Errorf("DisabledKatas = %v, want %v", cfg.DisabledKatas, want)
	}
}

func TestParseInlineList(t *testing.T) {
	cfg, err := Parse([]byte("disabled_katas: [ZC1, 'ZC2', \"ZC3\"]\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := []string{"ZC1", "ZC2", "ZC3"}; !reflect.DeepEqual(cfg.DisabledKatas, want) {
		t.Errorf("DisabledKatas = %v, want %v", cfg.DisabledKatas, want)
	}
}

func TestParseEmptyInlineList(t *testing.T) {
	cfg, err := Parse([]byte("disabled_katas: []\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.DisabledKatas) != 0 {
		t.Errorf("want empty, got %v", cfg.DisabledKatas)
	}
}

func TestParseBareListValue(t *testing.T) {
	cfg, err := Parse([]byte("disabled_katas: ZC9\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := []string{"ZC9"}; !reflect.DeepEqual(cfg.DisabledKatas, want) {
		t.Errorf("DisabledKatas = %v, want %v", cfg.DisabledKatas, want)
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
	if want := []string{"a", "b", "c", "d", "e", "f", "g", "h"}; !reflect.DeepEqual(got, want) {
		t.Errorf("colors = %v, want %v", got, want)
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

func TestParseEscapeEdgeCases(t *testing.T) {
	// Unknown escape kept verbatim; invalid \x falls back; trailing
	// backslash and lone NUL.
	cases := map[string]string{
		`"a\qb"`:    `a\qb`,
		`"\xZZ"`:    `\xZZ`,
		`"n\nl"`:    "n\nl",
		`"c\\d"`:    `c\d`,
		`"q\"q"`:    `q"q`,
		"\"r\\r0\"": "r\r0",
	}
	for in, want := range cases {
		cfg, err := Parse([]byte("error_color: " + in + "\n"))
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", in, err)
		}
		if cfg.ErrorColor != want {
			t.Errorf("%s => %q, want %q", in, cfg.ErrorColor, want)
		}
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
		":this is not yaml\n", // empty key
		"  - orphan\n",        // sequence item with no list key
		"no_color: maybe\n",   // invalid boolean
		"verbose: notabool\n", // invalid boolean
		"bareword\n",          // no key:value separator
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
	if !reflect.DeepEqual(cfg, Config{}) {
		t.Errorf("blank config should be zero value, got %+v", cfg)
	}
}
