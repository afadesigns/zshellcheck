// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1645(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — source os-release",
			input:    `source /etc/os-release`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — lsb_release -rs",
			input: `lsb_release -rs`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1645",
					Message: "`lsb_release` needs an optional package. Use `source /etc/os-release` and read `$ID` / `$VERSION_ID` / `$PRETTY_NAME` instead — always present, no fork.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — lsb_release -a",
			input: `lsb_release -a`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1645",
					Message: "`lsb_release` needs an optional package. Use `source /etc/os-release` and read `$ID` / `$VERSION_ID` / `$PRETTY_NAME` instead — always present, no fork.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1645")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
