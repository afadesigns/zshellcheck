// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package main

import (
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
