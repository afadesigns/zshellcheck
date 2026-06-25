// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

// TestZC1075FlagLedExpansionNotFlagged pins the flag-led-expansion false
// positive. A `${(%):-default}` / `${(P)…}` carries a parameter flag and a
// default the parser cannot fully model; it never produces an empty word,
// so the elision warning does not apply.
func TestZC1075FlagLedExpansionNotFlagged(t *testing.T) {
	for _, src := range []string{
		"echo ${(%):-default}",
		"echo ${(P):-x}",
		`read -q ${(%):-"?prompt "}`,
	} {
		if n := len(testutil.Check(src, "ZC1075")); n != 0 {
			t.Errorf("ZC1075 should not flag a flag-led expansion: %q (got %d)", src, n)
		}
	}
	// Ordinary unquoted scalars and array elements still elide and fire.
	for _, src := range []string{"echo $plain", "echo ${arr[1]}"} {
		if n := len(testutil.Check(src, "ZC1075")); n == 0 {
			t.Errorf("ZC1075 should still fire on %q", src)
		}
	}
}
