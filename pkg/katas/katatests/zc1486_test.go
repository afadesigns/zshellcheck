// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1486(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — curl https://host",
			input:    `curl https://host`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — curl -1 (TLSv1+)",
			input:    `curl -1 https://host`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — curl -2 (SSLv2)",
			input: `curl -2 https://host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1486",
					Message: "`curl -2` forces SSLv2/SSLv3 — removed from modern TLS libraries and subject to POODLE. Fix the server instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — curl -3 (SSLv3)",
			input: `curl -3 https://host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1486",
					Message: "`curl -3` forces SSLv2/SSLv3 — removed from modern TLS libraries and subject to POODLE. Fix the server instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1486")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
