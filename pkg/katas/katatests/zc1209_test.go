// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1209(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid systemctl with --no-pager",
			input:    `systemctl --no-pager status nginx`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid systemctl start",
			input:    `systemctl start nginx`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid systemctl status",
			input: `systemctl status nginx`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1209",
					Message: "Use `systemctl --no-pager` in scripts. Without it, systemctl invokes a pager that hangs in non-interactive execution.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1209")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
