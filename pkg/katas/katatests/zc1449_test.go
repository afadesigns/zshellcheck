// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1449(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — dnf install -y",
			input:    `dnf install -y vim`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — dnf install no -y",
			input: `dnf install vim`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1449",
					Message: "`dnf` without `-y` hangs on confirmation. Add `-y` for unattended runs.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — yum install no -y",
			input: `yum install httpd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1449",
					Message: "`yum` without `-y` hangs on confirmation. Add `-y` for unattended runs.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1449")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
