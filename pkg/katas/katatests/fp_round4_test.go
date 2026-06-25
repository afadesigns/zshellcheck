// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

// TestZC1098EvalAllArgsNotFlagged pins the argument-dispatch false
// positive. `eval "$@"` / `eval "$*"` run the positional parameters as
// the command (the standard async-worker idiom); quoting them with `(q)`
// collapses the words into one literal and breaks execution, exactly like
// the already-skipped `eval "$cmd"`.
func TestZC1098EvalAllArgsNotFlagged(t *testing.T) {
	for _, src := range []string{`eval "$@"`, `eval "$*"`} {
		if n := len(testutil.Check(src, "ZC1098")); n != 0 {
			t.Errorf("ZC1098 should not fire on the eval arg-dispatch idiom: %q (got %d)", src, n)
		}
	}
	// A variable used as data inside a command string still warrants (q).
	for _, src := range []string{`eval "echo $userdata"`, `eval "rm $path"`} {
		if n := len(testutil.Check(src, "ZC1098")); n == 0 {
			t.Errorf("ZC1098 should still fire on a variable embedded in an eval command: %q", src)
		}
	}
}
