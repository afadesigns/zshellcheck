// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1395(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — wait $pid",
			input:    `wait $pid`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — wait -n",
			input: `wait -n`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1395",
					Message: "`wait -n` is Bash 4.3+. Zsh's `wait` waits on specific PIDs/jobs or (bare `wait`) all jobs. For any-child semantics, loop over PIDs with individual `wait $pid` calls.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1395")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
