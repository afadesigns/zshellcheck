// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

// TestZC1043DeclaredLocalsNotFlagged pins the declared-local false
// positive. A flagged `declare`/`typeset` parses as a DeclarationStatement
// rather than a command, so its names must still be recognised as local;
// a later assignment to one of them is not an accidental global.
func TestZC1043DeclaredLocalsNotFlagged(t *testing.T) {
	for _, src := range []string{
		"f() {\n  declare -a versions\n  versions=(a b)\n}",
		"f() {\n  typeset -i count\n  count=5\n}",
	} {
		if n := len(testutil.Check(src, "ZC1043")); n != 0 {
			t.Errorf("ZC1043 should not flag an assignment to a declared local: %q (got %d)", src, n)
		}
	}
	if n := len(testutil.Check("f() { total=5 }", "ZC1043")); n == 0 {
		t.Error("ZC1043 should still fire on a bare unscoped assignment")
	}
}

// TestZC1043SpecialTiedParamsNotFlagged pins the special-parameter false
// positive. path/fpath/cdpath and the other tied or shell-state
// parameters are global by nature; `local` scopes them away and breaks
// plugin load/unload, so ZC1043 must not suggest it.
func TestZC1043SpecialTiedParamsNotFlagged(t *testing.T) {
	for _, src := range []string{
		"f() { fpath=(/x $fpath) }",
		"f() { path=(/x $path) }",
		"f() { cdpath=(/a /b) }",
		"f() { manpath=(/m) }",
	} {
		if n := len(testutil.Check(src, "ZC1043")); n != 0 {
			t.Errorf("ZC1043 should not flag a special tied parameter: %q (got %d)", src, n)
		}
	}
}
