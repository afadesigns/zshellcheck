// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1579(t *testing.T) {
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
			name:     "valid — curl --retry-all-errors --max-time 30 URL",
			input:    `curl https://host --retry-all-errors --max-time 30`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — curl --retry-all-errors -m 30 URL",
			input:    `curl https://host --retry-all-errors -m 30`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — curl URL --retry-all-errors (no max-time)",
			input: `curl https://host --retry-all-errors`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1579",
					Message: "`curl --retry-all-errors` with no `--max-time` hammers the upstream on failure. Pair with `-m <seconds>` or use `--retry-connrefused`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1579")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
