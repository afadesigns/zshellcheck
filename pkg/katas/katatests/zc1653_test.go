// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1653(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — $sysparams[pid]",
			input:    `echo "$sysparams[pid]"`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — plain PID reference",
			input:    `echo "$PPID"`,
			expected: []katas.Violation{},
		},
		{
			name:  `invalid — echo "$BASHPID"`,
			input: `echo "$BASHPID"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1653",
					Message: "`$BASHPID` is Bash-only. Use `$sysparams[pid]` after `zmodload zsh/system`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  `invalid — print "sub=${BASHPID}"`,
			input: `print -r -- "sub=${BASHPID}"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1653",
					Message: "`$BASHPID` is Bash-only. Use `$sysparams[pid]` after `zmodload zsh/system`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1653")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
