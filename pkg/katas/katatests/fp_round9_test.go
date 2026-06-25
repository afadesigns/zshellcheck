// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

// TestZC1075WidthAndTypesetArray pins two more array/elision refinements.
func TestZC1075WidthAndTypesetArray(t *testing.T) {
	skip := []string{
		"cmd ${(l:5:)x}",              // left-pad — fixed width, never empty
		"cmd ${(r:3:)y}",              // right-pad
		"typeset arr=(a b)\nuse $arr", // unflagged typeset array literal
	}
	for _, src := range skip {
		if n := len(testutil.Check(src, "ZC1075")); n != 0 {
			t.Errorf("ZC1075 should not flag %q (got %d)", src, n)
		}
	}
	fire := []string{
		"print ${(r)things}",              // reverse without width can be empty
		"typeset q=\"(literal)\"\nuse $q", // quoted scalar, not an array
	}
	for _, src := range fire {
		if n := len(testutil.Check(src, "ZC1075")); n == 0 {
			t.Errorf("ZC1075 should still fire on %q", src)
		}
	}
}
