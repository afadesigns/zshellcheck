// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

// TestZC1097PositionalLoopVar pins the positional-loop-variable false
// positive. `for 1 in …` binds the positional `$1`, which is already
// scoped to the function and cannot take `local` — `local 1` is a zsh
// error — so the kata must not suggest it.
func TestZC1097PositionalLoopVar(t *testing.T) {
	for _, src := range []string{
		"f() { for 1 in $p; do :; done }",
		"f() { for 2 in a b; do :; done }",
	} {
		if n := len(testutil.Check(src, "ZC1097")); n != 0 {
			t.Errorf("ZC1097 should not fire on a positional loop var: %q (got %d)", src, n)
		}
	}
	if n := len(testutil.Check("f() { for x in $p; do :; done }", "ZC1097")); n == 0 {
		t.Error("ZC1097 should still fire on a named loop var")
	}
}

// TestZC1049ListingVsDefinition pins the alias-listing false positive.
// The listing and query forms print existing aliases and have no
// name=value, so no function can replace them; only real definitions
// should fire.
func TestZC1049ListingVsDefinition(t *testing.T) {
	for _, src := range []string{
		"alias",
		"alias -L",
		"alias $1",
		"x=$(alias | cut -d '=' -f 1)",
	} {
		if n := len(testutil.Check(src, "ZC1049")); n != 0 {
			t.Errorf("ZC1049 should not fire on an alias listing form: %q (got %d)", src, n)
		}
	}
	for _, src := range []string{
		"alias d='dirs -v'",
		"alias -- -='cd -'",
		"alias ${ZSHZ_CMD:-z}='zshz 2>&1'",
	} {
		if n := len(testutil.Check(src, "ZC1049")); n == 0 {
			t.Errorf("ZC1049 should fire on an alias definition: %q", src)
		}
	}
}
