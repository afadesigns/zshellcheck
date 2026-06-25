// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

// TestZC1044CdAssignmentNotFlagged pins the `cd=value` false positive. An
// assignment to a variable named `cd` on the right of `&&`/`||` is parsed
// as a `cd` command whose first argument is the `=value` tail; it is not a
// directory change, so it must not be flagged.
func TestZC1044CdAssignmentNotFlagged(t *testing.T) {
	for _, src := range []string{
		"foo && cd=$bar",
		"(( ZSHZ_TILDE )) && cd=${cd/#x/y}",
	} {
		if n := len(testutil.Check(src, "ZC1044")); n != 0 {
			t.Errorf("ZC1044 should not flag an assignment to a var named cd: %q (got %d)", src, n)
		}
	}
	for _, src := range []string{"cd /tmp", "foo && cd /tmp"} {
		if n := len(testutil.Check(src, "ZC1044")); n == 0 {
			t.Errorf("ZC1044 should still flag an unguarded cd: %q", src)
		}
	}
}
