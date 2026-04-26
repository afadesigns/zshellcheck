// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1736(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `pulumi preview`",
			input:    `pulumi preview`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `pulumi up` (interactive prompt kept)",
			input:    `pulumi up`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `pulumi stack ls --yes` (not up/destroy/refresh)",
			input:    `pulumi stack ls --yes`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `pulumi destroy --yes`",
			input: `pulumi destroy --yes`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1736",
					Message: "`pulumi destroy --yes` skips the preview-and-confirm — a misresolved stack or credential wipes / mutates infrastructure with no review. Gate behind `pulumi preview` plus a manual approval step.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `pulumi up -y`",
			input: `pulumi up -y`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1736",
					Message: "`pulumi up -y` skips the preview-and-confirm — a misresolved stack or credential wipes / mutates infrastructure with no review. Gate behind `pulumi preview` plus a manual approval step.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `pulumi refresh --yes`",
			input: `pulumi refresh --yes`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1736",
					Message: "`pulumi refresh --yes` skips the preview-and-confirm — a misresolved stack or credential wipes / mutates infrastructure with no review. Gate behind `pulumi preview` plus a manual approval step.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1736")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
