// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1994(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `lvreduce -L 10G $LV` (interactive confirm)",
			input:    `lvreduce -L 10G $LV`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `lvextend -L +10G $LV` (grow)",
			input:    `lvextend -L +10G $LV`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `lvreduce -f -L 10G $LV`",
			input: `lvreduce -f -L 10G $LV`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1994",
					Message: "`lvreduce -f` skips the shrink-confirmation prompt — the filesystem above still believes the tail is allocated and the next mount sees corruption. Shrink fs first, or use `lvreduce --resizefs`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `lvreduce -y -L 10G $LV`",
			input: `lvreduce -y -L 10G $LV`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1994",
					Message: "`lvreduce -y` skips the shrink-confirmation prompt — the filesystem above still believes the tail is allocated and the next mount sees corruption. Shrink fs first, or use `lvreduce --resizefs`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1994")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
