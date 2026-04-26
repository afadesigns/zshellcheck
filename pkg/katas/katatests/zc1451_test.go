// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1451(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — pip install --user",
			input:    `pip install --user requests`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — pip install (system-wide)",
			input: `pip install requests`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1451",
					Message: "`pip install` without `--user` or an active venv targets system Python. Use `python -m venv` / `uv` / `--user` for scoped installs.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1451")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
