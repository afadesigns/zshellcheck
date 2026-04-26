// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1232(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid pip install --user",
			input:    `pip install --user requests`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid pip3 list",
			input:    `pip3 list`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid bare pip install",
			input: `pip install requests`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1232",
					Message: "Use `pip install --user` or a virtualenv instead of bare `pip install`. System-wide pip installs can break OS package managers.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1232")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
