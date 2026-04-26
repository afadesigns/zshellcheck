// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package main

import (
	"bytes"
	"flag"
	"strings"
	"testing"
)

func TestWrapShort(t *testing.T) {
	got := wrap("short line", 80)
	if len(got) != 1 || got[0] != "short line" {
		t.Errorf("short line wrapped unexpectedly: %#v", got)
	}
}

func TestWrapBreaks(t *testing.T) {
	got := wrap("alpha beta gamma delta epsilon", 12)
	if len(got) < 2 {
		t.Errorf("expected wrap to break, got: %#v", got)
	}
	for _, line := range got {
		if len(line) > 12 {
			t.Errorf("wrap exceeded width: %q", line)
		}
	}
}

func TestFlagValueType(t *testing.T) {
	fs := flag.NewFlagSet("t", flag.ContinueOnError)
	fs.String("s", "x", "")
	fs.Int("i", 0, "")
	fs.Bool("b", false, "")
	fs.Float64("f", 0, "")

	cases := map[string]string{
		"s": "string",
		"i": "int",
		"b": "",
		"f": "float",
	}
	for name, want := range cases {
		f := fs.Lookup(name)
		if got := flagValueType(f); got != want {
			t.Errorf("flagValueType(%s) = %q, want %q", name, got, want)
		}
	}
}

func TestPrintUsageContainsCoreSections(t *testing.T) {
	fs := flag.NewFlagSet("zshellcheck", flag.ContinueOnError)
	fs.String("format", "text", "Output format")
	fs.String("severity", "", "Severity filter")
	fs.Bool("fix", false, "Apply auto-fixes")
	fs.Bool("diff", false, "Preview diff")
	fs.Bool("dry-run", false, "Dry-run")
	fs.String("cpuprofile", "", "CPU profile path")
	fs.Bool("version", false, "Print version")
	fs.Bool("no-color", false, "Disable colour")
	fs.Bool("no-banner", false, "Suppress banner")
	fs.Bool("verbose", false, "Verbose")

	var buf bytes.Buffer
	printUsage(&buf, fs, false)
	out := buf.String()

	for _, want := range []string{"USAGE", "OUTPUT", "FILTER", "AUTO-FIX", "DIAGNOSTICS", "EXAMPLES", "DOCUMENTATION"} {
		if !strings.Contains(out, want) {
			t.Errorf("printUsage missing %q section", want)
		}
	}
}

func TestRenderFlag(t *testing.T) {
	fs := flag.NewFlagSet("t", flag.ContinueOnError)
	fs.String("config", "default.yml", "Path to the config file")

	var buf bytes.Buffer
	renderFlag(&buf, palette{}, fs.Lookup("config"))

	out := buf.String()
	if !strings.Contains(out, "-config") {
		t.Errorf("renderFlag missing flag name: %q", out)
	}
	if !strings.Contains(out, "Path to the config file") {
		t.Errorf("renderFlag missing usage text: %q", out)
	}
}

func TestPaletteDisabled(t *testing.T) {
	p := palette{enabled: false}
	if got := p.bold("x"); got != "x" {
		t.Errorf("disabled bold returns wrapped: %q", got)
	}
	if got := p.dim("x"); got != "x" {
		t.Errorf("disabled dim returns wrapped: %q", got)
	}
	if got := p.section("x"); got != "x" {
		t.Errorf("disabled section returns wrapped: %q", got)
	}
}

func TestPaletteEnabled(t *testing.T) {
	p := palette{enabled: true}
	out := p.bold("x")
	if !strings.Contains(out, "\x1b[1m") || !strings.Contains(out, "\x1b[0m") {
		t.Errorf("enabled bold missing ANSI: %q", out)
	}
}

func TestPaletteAllAccents(t *testing.T) {
	p := palette{enabled: true}
	for name, fn := range map[string]func(string) string{
		"dim":      p.dim,
		"section":  p.section,
		"flagName": p.flagName,
		"link":     p.link,
	} {
		if got := fn("x"); !strings.Contains(got, "\x1b[") {
			t.Errorf("%s did not emit ANSI escape: %q", name, got)
		}
	}
}
