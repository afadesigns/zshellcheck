// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1408(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — echo $FPATH",
			input:    `echo $FPATH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $BASH_FUNC_myfn",
			input: `echo $BASH_FUNC_myfn`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1408",
					Message: "`BASH_FUNC_*` exported-function envvars are Bash-only. Zsh does not consume them; export function definitions via `autoload` + `$FPATH` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1408")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
