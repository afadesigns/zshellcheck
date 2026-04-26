// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1520(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — read varname",
			input:    `read varname`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — vared myvar",
			input: `vared myvar`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1520",
					Message: "`vared` requires a TTY — in a non-interactive script it errors or hangs. Use `read`, stdin, or environment variables for scripted input.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1520")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
