// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1245(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid curl with TLS",
			input:    `curl -fsSL https://example.com`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid curl -k",
			input: `curl -k https://example.com`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1245",
					Message: "Avoid `curl -k`/`--insecure` — it disables TLS certificate verification. Fix the certificate chain or use `--cacert` to specify a CA bundle.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1245")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
