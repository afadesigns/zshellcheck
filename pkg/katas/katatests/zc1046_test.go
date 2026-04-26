// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1046(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "direct eval",
			input: `eval "echo hello"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1046",
					Message: "Avoid `eval`. It allows execution of arbitrary code and is hard to debug.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "not eval",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:  "builtin eval",
			input: `builtin eval "echo hello"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1046",
					Message: "Avoid `eval`. It allows execution of arbitrary code and is hard to debug.",
					Line:    1,
					Column:  9,
				},
			},
		},
		{
			name:  "command eval",
			input: `command eval "echo hello"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1046",
					Message: "Avoid `eval`. It allows execution of arbitrary code and is hard to debug.",
					Line:    1,
					Column:  9,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1046")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
