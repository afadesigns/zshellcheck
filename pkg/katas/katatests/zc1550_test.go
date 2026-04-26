// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1550(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — apt-mark unhold pkg",
			input:    `apt-mark unhold openssh-server`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — apt-mark showhold",
			input:    `apt-mark showhold`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — apt-mark hold pkg",
			input: `apt-mark hold openssh-server`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1550",
					Message: "`apt-mark hold` pins the package — blocks future CVE fixes. Document the reason and schedule an unhold review.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — apt-mark hold multiple pkgs",
			input: `apt-mark hold openssh-server libc6`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1550",
					Message: "`apt-mark hold` pins the package — blocks future CVE fixes. Document the reason and schedule an unhold review.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1550")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
