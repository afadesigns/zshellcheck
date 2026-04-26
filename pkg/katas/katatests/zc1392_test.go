// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1392(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — unrelated echo",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $CHILD_MAX",
			input: `echo $CHILD_MAX`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1392",
					Message: "`$CHILD_MAX` is Bash-only. Zsh uses `limit -s maxproc` or `ulimit -u` for process-count limits.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1392")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
