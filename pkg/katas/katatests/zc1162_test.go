// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1162(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid cp -a",
			input:    `cp -a src/ dst/`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid cp single file",
			input:    `cp file.txt backup.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid cp -r",
			input: `cp -r src/ dst/`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1162",
					Message: "Use `cp -a` instead of `cp -r` to preserve permissions, timestamps, and symlinks. Archive mode ensures a faithful copy.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1162")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
