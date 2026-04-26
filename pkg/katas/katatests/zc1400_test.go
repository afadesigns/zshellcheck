// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1400(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — awk for other field",
			input:    `awk '{print $1}' file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — cut on HOSTTYPE",
			input: `cut -d- -f1 $HOSTTYPE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1400",
					Message: "Use Zsh `$CPUTYPE` for pure architecture instead of splitting `$HOSTTYPE` with `cut`/`awk`/`sed`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1400")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
