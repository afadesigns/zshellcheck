// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package fix

import (
	"strings"
	"testing"
)

// These regressions pin three auto-fixes that previously emitted broken
// or wrong source. Each asserts the fix no longer corrupts its input.
// They were found by auditing real-corpus findings against runnable Zsh.

func TestFixRegression_ZC1086_NamespacedFunctionName(t *testing.T) {
	// `function name"${1:-}"() { ... }` builds the function name from a
	// parameter expansion (gitstatus does this for per-instance
	// namespacing). The fix used to insert `()` after the parsed
	// identifier, splitting the real name so a different function was
	// defined. It must now leave the source untouched.
	src := "function gitstatus_query\"${1:-}\"() { print body; }\n"
	if got := runFix(t, src); got != src {
		t.Errorf("ZC1086 rewrote a parameter-expanded function name\ninput: %q\ngot:   %q", src, got)
	}
}

func TestFixRegression_ZC1010_BinaryDashODashA(t *testing.T) {
	// `[ x -o y ]` is a POSIX OR and `[ a -a b ]` an AND. `[[ ]]` has no
	// binary `-o`/`-a` (it uses `||`/`&&` and reads them as unary tests),
	// so a naive `[` -> `[[` swap is a syntax error. The fix must bail.
	for _, src := range []string{
		"[ \"$a\" = \"\" -o \"$b\" = \"\" ] && print y\n",
		"[ -f a -a -f b ]\n",
	} {
		if got := runFix(t, src); strings.Contains(got, "[[") {
			t.Errorf("ZC1010 converted a binary -a/-o test to invalid [[ ]]\ninput: %q\ngot:   %q", src, got)
		}
	}
}

func TestFixRegression_ZC1010_UnaryStillConverts(t *testing.T) {
	// The unary `[ -o opt ]` (option test) stays fixable; only binary
	// `-a`/`-o` bail. This guards the bail from over-reaching.
	if got := runFix(t, "[ -o nounset ]\n"); !strings.Contains(got, "[[ -o nounset ]]") {
		t.Errorf("ZC1010 should still convert a unary -o test, got %q", got)
	}
}

func TestFixRegression_ZC1076_NoDuplicateFlags(t *testing.T) {
	// The fix used to insert ` -Uz` unconditionally, so `autoload -U …`
	// became `autoload -Uz -U …`. It must add only the absent flag.
	got := runFix(t, "autoload -U is-at-least\n")
	if strings.Contains(got, "-Uz -U") || strings.Contains(got, "-U -U") {
		t.Errorf("ZC1076 produced duplicate flags: %q", got)
	}
	if !strings.Contains(got, "-z") {
		t.Errorf("ZC1076 should add the missing -z flag, got %q", got)
	}
}
