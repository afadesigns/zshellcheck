// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1339(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid wc -l with file",
			input:    `wc -l file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid wc -c",
			input:    `wc -c`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid wc -l in pipeline",
			input: `wc -l`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1339",
					Message: "Use Zsh `${#${(f)var}}` for line counting instead of piping through `wc -l`. Parameter expansion avoids spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1339")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
