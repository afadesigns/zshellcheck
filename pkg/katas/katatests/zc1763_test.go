// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1763(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `docker compose down`",
			input:    `docker compose down`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `docker compose up -d`",
			input:    `docker compose up -d`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `docker compose down -v`",
			input: `docker compose down -v`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1763",
					Message: "`docker compose down -v` wipes every named volume declared in the stack — database, cache, uploaded assets go with it. Drop the flag in CI / prod scripts.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `docker compose down --volumes`",
			input: `docker compose down --volumes`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1763",
					Message: "`docker compose down --volumes` wipes every named volume declared in the stack — database, cache, uploaded assets go with it. Drop the flag in CI / prod scripts.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `docker-compose down -v` (hyphen form)",
			input: `docker-compose down -v`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1763",
					Message: "`docker compose down -v` wipes every named volume declared in the stack — database, cache, uploaded assets go with it. Drop the flag in CI / prod scripts.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1763")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
