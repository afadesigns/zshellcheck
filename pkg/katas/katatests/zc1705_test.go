// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1705(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — awk without -i",
			input:    `awk '{print}' file`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — awk -i with non-inplace path",
			input:    `awk -i /usr/share/awk/lib.awk '{print}' file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — awk -i inplace",
			input: `awk -i inplace '{print}' file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1705",
					Message: "`awk -i inplace` is gawk-only — fails on mawk / BSD awk / busybox awk. For portability rewrite as `awk … input > tmp && mv tmp input`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — gawk -i inplace -v",
			input: `gawk -i inplace -v x=1 '{print x}' file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1705",
					Message: "`awk -i inplace` is gawk-only — fails on mawk / BSD awk / busybox awk. For portability rewrite as `awk … input > tmp && mv tmp input`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1705")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
