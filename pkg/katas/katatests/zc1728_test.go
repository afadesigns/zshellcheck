// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1728(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `pip install --index-url https://pypi.org/simple pkg`",
			input:    `pip install --index-url https://pypi.org/simple pkg`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `pip install pkg` (default https index)",
			input:    `pip install pkg`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `pip install --index-url http://internal/simple pkg`",
			input: `pip install --index-url http://internal/simple pkg`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1728",
					Message: "`pip install --index-url http://internal/simple` fetches packages over plaintext HTTP — any MITM swaps the wheel for code execution on the host. Use `https://`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `pip install -i http://internal/simple pkg`",
			input: `pip install -i http://internal/simple pkg`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1728",
					Message: "`pip install --index-url http://internal/simple` fetches packages over plaintext HTTP — any MITM swaps the wheel for code execution on the host. Use `https://`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `pip install --extra-index-url=http://internal/simple pkg`",
			input: `pip install --extra-index-url=http://internal/simple pkg`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1728",
					Message: "`pip install --index-url --extra-index-url=http://internal/simple` fetches packages over plaintext HTTP — any MITM swaps the wheel for code execution on the host. Use `https://`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1728")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
