// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1689(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — borg delete without --force (prompts)",
			input:    `borg delete /backup::archive1`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — borg prune",
			input:    `borg prune --keep-last 7 /backup`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — borg delete --force",
			input: `borg delete --force /backup`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1689",
					Message: "`borg delete --force` skips confirmation and can nuke the whole repository on a typo — use `borg prune --keep-*` with a retention policy, or gate outright deletion behind a manual review.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1689")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
