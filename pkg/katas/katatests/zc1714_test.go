// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1714(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — gh repo delete without --yes (prompts)",
			input:    `gh repo delete owner/repo`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — gh repo create",
			input:    `gh repo create owner/repo`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — gh repo delete --yes",
			input: `gh repo delete owner/repo --yes`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1714",
					Message: "`gh repo delete --yes` bypasses GitHub's confirmation — a typo or stale variable destroys the target with no soft-delete. Drop `--yes` so a human confirms.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — gh release delete --yes",
			input: `gh release delete v1.0 --yes`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1714",
					Message: "`gh release delete --yes` bypasses GitHub's confirmation — a typo or stale variable destroys the target with no soft-delete. Drop `--yes` so a human confirms.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1714")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
