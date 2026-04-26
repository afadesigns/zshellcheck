// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1573(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — chattr =i (set exclusive)",
			input:    `chattr =i /etc/shadow`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — chattr -i /etc/shadow",
			input: `chattr -i /etc/shadow`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1573",
					Message: "`chattr -i` removes the tamper-evident attribute. If this is a one-shot upgrade, re-set the attribute at the end of the block.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — chattr -a /var/log/auth.log",
			input: `chattr -a /var/log/auth.log`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1573",
					Message: "`chattr -a` removes the tamper-evident attribute. If this is a one-shot upgrade, re-set the attribute at the end of the block.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1573")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
