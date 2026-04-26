// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1718(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `gh secret set NAME --body-file path`",
			input:    `gh secret set NAME --body-file /run/secrets/foo`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `gh secret set NAME --body -` (read stdin)",
			input:    `gh secret set NAME --body -`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `gh variable set NAME --body val` (non-secret)",
			input:    `gh variable set NAME --body publicvalue`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `gh secret set NAME --body SECRET`",
			input: `gh secret set NAME --body hunter2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1718",
					Message: "`gh secret set ... --body hunter2` puts the secret in argv — visible in `ps`, `/proc`, history. Use `--body-file PATH` or pipe via stdin (`... --body -` with `printf %s \"$SECRET\" |`).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `gh secret set NAME --body=SECRET`",
			input: `gh secret set NAME --body=hunter2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1718",
					Message: "`gh secret set ... --body=hunter2` puts the secret in argv — visible in `ps`, `/proc`, history. Use `--body-file PATH` or pipe via stdin (`... --body -` with `printf %s \"$SECRET\" |`).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `gh secret set NAME -b SECRET`",
			input: `gh secret set NAME -b hunter2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1718",
					Message: "`gh secret set ... --body hunter2` puts the secret in argv — visible in `ps`, `/proc`, history. Use `--body-file PATH` or pipe via stdin (`... --body -` with `printf %s \"$SECRET\" |`).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1718")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
