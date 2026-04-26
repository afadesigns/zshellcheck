// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1428(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — curl with no auth",
			input:    `curl https://example.com`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — curl -u user:pass",
			input: `curl -u alice:secret123 https://example.com`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1428",
					Message: "`curl -u user:pass` leaks credentials into the process list. Use `-u user:` (prompt), `--netrc`, or a credentials manager.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1428")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
