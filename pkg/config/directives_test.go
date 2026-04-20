package config

import (
	"reflect"
	"testing"
)

func TestParseDirectives_Trailing(t *testing.T) {
	src := `echo hi  # zshellcheck disable=ZC1075
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

func TestParseDirectives_PrecedingOwnLine(t *testing.T) {
	src := `# zshellcheck disable=ZC1136,ZC1141
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

func TestParseDirectives_FileTail(t *testing.T) {
	// Directive at file end with no code after it becomes file-wide.
	src := `echo hi
# zshellcheck disable=ZC1075
`
	d := ParseDirectives(src)
	if !d.IsDisabledOn("ZC1075", 42) {
		t.Errorf("expected ZC1075 disabled file-wide")
	}
}

func TestParseDirectives_MultipleIDs(t *testing.T) {
	src := `rm -rf /tmp/x # zshellcheck disable=ZC1136, ZC1075
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
