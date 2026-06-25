// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

// TestCommandPositionAssignmentNotMisflagged pins the command-position
// assignment fix. A `name=value` after `&&` / `||` or in a pipeline is an
// assignment, not a command; a variable named for a builtin was previously
// parsed as that builtin and mis-flagged. The common guard idiom
// `[[ cond ]] && var=value` is where this surfaced across the corpus.
func TestCommandPositionAssignmentNotMisflagged(t *testing.T) {
	cases := []struct {
		src   string
		katas []string
	}{
		// `test`, `seq`, `timeout` as variable names, not commands.
		{`[[ -n $x ]] && test=$(foo)`, []string{"ZC1006", "ZC1020", "ZC1036", "ZC1293"}},
		{`[[ $a = b ]] && seq=$seq_mac`, []string{"ZC1061", "ZC1276"}},
		{`(( $2 >= 0 )) && timeout=-t$2`, []string{"ZC1167"}},
	}
	for _, tc := range cases {
		for _, kata := range tc.katas {
			if n := len(testutil.Check(tc.src, kata)); n != 0 {
				t.Errorf("%s wrongly fired on a command-position assignment: %q (got %d)", kata, tc.src, n)
			}
		}
	}
}

// TestCommandPositionAssignmentKeepsRealCommands confirms the reshape does
// not hide real commands: the builtin invoked for real after a logical
// operator still draws its finding, and an env-prefixed command keeps both
// the assignment and the command visible.
func TestCommandPositionAssignmentKeepsRealCommands(t *testing.T) {
	if n := len(testutil.Check(`x && timeout 5 cmd`, "ZC1167")); n != 1 {
		t.Errorf("a real `timeout` command after `&&` should still fire ZC1167, got %d", n)
	}
	// A genuine array append assignment in command position is now visible.
	if n := len(testutil.Check(`[[ -z $x ]] && fpath=($fpath $dir)`, "ZC1071")); n != 1 {
		t.Errorf("a real array append after `&&` should fire ZC1071, got %d", n)
	}
}
