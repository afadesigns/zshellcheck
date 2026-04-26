// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1666(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — kubectl patch --type=strategic",
			input:    `kubectl patch deployment nginx --type=strategic`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — kubectl patch --type=merge",
			input:    `kubectl patch deployment nginx --type=merge`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — kubectl patch --type=json joined",
			input: `kubectl patch deployment nginx --type=json -p '[...]'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1666",
					Message: "`kubectl patch --type=json` applies a raw RFC-6902 patch that bypasses strategic-merge reconciliation — prefer `--type=strategic` and hold JSON patches behind code review.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — kubectl patch --type json split",
			input: `kubectl patch deployment nginx --type json -p '[...]'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1666",
					Message: "`kubectl patch --type=json` applies a raw RFC-6902 patch that bypasses strategic-merge reconciliation — prefer `--type=strategic` and hold JSON patches behind code review.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1666")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
