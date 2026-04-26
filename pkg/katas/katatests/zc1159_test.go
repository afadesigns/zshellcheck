// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1159(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid tar czf",
			input:    `tar czf archive.tar.gz dir`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid tar extract",
			input:    `tar xf archive.tar.gz`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid tar cf without compression",
			input: `tar cf archive.tar dir`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1159",
					Message: "Specify an explicit compression flag (`-z`, `-j`, `-J`) when creating tar archives. Relying on auto-detection reduces clarity and portability.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1159")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
