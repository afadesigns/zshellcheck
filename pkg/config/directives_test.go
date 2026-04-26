// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package config

import (
	"reflect"
	"testing"
)

func TestParseDirectives_TrailingWithIDs(t *testing.T) {
	src := `echo hi  # noka: ZC1075
echo again
`
	d := ParseDirectives(src)
	if !d.IsDisabledOn("ZC1075", 1) {
		t.Errorf("expected ZC1075 disabled on line 1")
	}
	if d.IsDisabledOn("ZC1075", 2) {
		t.Errorf("did not expect ZC1075 disabled on line 2")
	}
}

func TestParseDirectives_TrailingBareSilencesAll(t *testing.T) {
	src := `rm -rf $target  # noka
echo after
`
	d := ParseDirectives(src)
	if !d.IsDisabledOn("ZC1136", 1) || !d.IsDisabledOn("ZC9999", 1) {
		t.Errorf("expected bare # noka to silence every kata on line 1")
	}
	if d.IsDisabledOn("ZC1136", 2) {
		t.Errorf("did not expect # noka to leak onto line 2")
	}
}

func TestParseDirectives_PrecedingWithIDs(t *testing.T) {
	src := `# noka: ZC1136, ZC1141
rm -rf /tmp/noisy
echo after
`
	d := ParseDirectives(src)
	if !d.IsDisabledOn("ZC1136", 2) {
		t.Errorf("expected ZC1136 disabled on line 2")
	}
	if !d.IsDisabledOn("ZC1141", 2) {
		t.Errorf("expected ZC1141 disabled on line 2")
	}
	if d.IsDisabledOn("ZC1136", 3) {
		t.Errorf("did not expect ZC1136 disabled on line 3")
	}
}

func TestParseDirectives_PrecedingBareSilencesAll(t *testing.T) {
	src := `# noka
rm -rf /tmp/noisy
echo after
`
	d := ParseDirectives(src)
	if !d.IsDisabledOn("ZC1136", 2) || !d.IsDisabledOn("ZC9999", 2) {
		t.Errorf("expected bare # noka to silence every kata on the next code line")
	}
	if d.IsDisabledOn("ZC1136", 3) {
		t.Errorf("did not expect bare # noka to leak past the next code line")
	}
}

func TestParseDirectives_FileTail(t *testing.T) {
	// Directive at file end with no code after it becomes file-wide.
	src := `echo hi
# noka: ZC1075
`
	d := ParseDirectives(src)
	if !d.IsDisabledOn("ZC1075", 42) {
		t.Errorf("expected ZC1075 disabled file-wide")
	}
}

func TestParseDirectives_FileTailBareSilencesAll(t *testing.T) {
	src := `echo hi
# noka
`
	d := ParseDirectives(src)
	if !d.FileAll {
		t.Errorf("expected FileAll set by file-tail bare # noka")
	}
	if !d.IsDisabledOn("ZC1234", 42) || !d.IsDisabledOn("ZC9999", 1) {
		t.Errorf("expected file-wide silence to apply across every line")
	}
}

func TestParseDirectives_MultipleIDs(t *testing.T) {
	src := `rm -rf /tmp/x # noka: ZC1136, ZC1075
`
	d := ParseDirectives(src)
	if !reflect.DeepEqual(d.PerLine[1], []string{"ZC1136", "ZC1075"}) {
		t.Errorf("expected [ZC1136 ZC1075] on line 1, got %v", d.PerLine[1])
	}
}

func TestParseDirectives_None(t *testing.T) {
	d := ParseDirectives("echo hello\n")
	if d.HasAny() {
		t.Error("expected no directives")
	}
}

func TestParseDirectives_NokaInsideWordIgnored(t *testing.T) {
	// `noka` must match as a whole word — substring inside another token
	// should not be misread as a directive.
	src := `echo nokaroni  # nokarama is not a directive
`
	d := ParseDirectives(src)
	if d.HasAny() {
		t.Errorf("expected no directives when 'noka' appears as a substring of another word, got %+v", d)
	}
}

func TestParseDirectives_LegacyZshellcheckFormDoesNothing(t *testing.T) {
	// The old `# zshellcheck disable=…` syntax was retired in v1.0.15. Make
	// sure it now silently does nothing — the test doubles as a regression
	// guard against accidental re-introduction.
	src := `rm -rf $target  # zshellcheck disable=ZC1136
`
	d := ParseDirectives(src)
	if d.HasAny() {
		t.Errorf("legacy directive should not register, got %+v", d)
	}
}
