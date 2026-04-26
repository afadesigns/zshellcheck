// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1465(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — setenforce 1",
			input:    `setenforce 1`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — setenforce Enforcing",
			input:    `setenforce Enforcing`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — setenforce 0",
			input: `setenforce 0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1465",
					Message: "`setenforce 0` disables SELinux enforcement host-wide. Fix the AVC with `audit2allow` instead and keep enforcing mode on.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — setenforce Permissive",
			input: `setenforce Permissive`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1465",
					Message: "`setenforce 0` disables SELinux enforcement host-wide. Fix the AVC with `audit2allow` instead and keep enforcing mode on.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1465")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
