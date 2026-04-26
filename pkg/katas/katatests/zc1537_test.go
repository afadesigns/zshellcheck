// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1537(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — lvremove vg/lv",
			input:    `lvremove vg0/lv0`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — lvremove -f vg/lv",
			input: `lvremove -f vg0/lv0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1537",
					Message: "`lvremove -f` skips the confirmation — a typo in the volume name destroys every filesystem on top of it.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — vgremove -f vg",
			input: `vgremove -f vg0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1537",
					Message: "`vgremove -f` skips the confirmation — a typo in the volume name destroys every filesystem on top of it.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — pvremove -ff pv",
			input: `pvremove -ff /devicenode`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1537",
					Message: "`pvremove -ff` skips the confirmation — a typo in the volume name destroys every filesystem on top of it.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1537")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
