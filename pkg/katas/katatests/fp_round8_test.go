// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

// TestZC1075ArrayNotFlagged pins the array-quoting false positive. A bare
// expansion of a variable declared as an array joins its elements into one
// word under quoting, so the elision advice (`"$var"`) is wrong; the
// correct form is `"${arr[@]}"`. ZC1075 now skips known arrays.
func TestZC1075ArrayNotFlagged(t *testing.T) {
	arrays := []string{
		"local -a flags\njj $flags",
		"flags=(--foo --bar)\ngit $flags",
		"typeset -A m\nprint $m",
		"local opts=(x y)\ncmd $opts",
		"declare -a items\nrun $items",
		"readonly -a paths\nuse $paths",
	}
	for _, src := range arrays {
		if n := len(testutil.Check(src, "ZC1075")); n != 0 {
			t.Errorf("ZC1075 should not flag a known array: %q (got %d)", src, n)
		}
	}
	// A scalar must not be mistaken for an array. The renderer prints a
	// scalar value `$1` as `name=($1)`, so detection reads the parsed value
	// node, not the text; these must all still fire.
	scalars := []string{
		"local worker=$1\nprint $worker",
		"local x=$(cmd)\necho $x",
		"local p=$a/$b\nrm $p",
		"typeset -i n\necho $n",
		"rm $file",
	}
	for _, src := range scalars {
		if n := len(testutil.Check(src, "ZC1075")); n == 0 {
			t.Errorf("ZC1075 should still flag a scalar: %q", src)
		}
	}
}
