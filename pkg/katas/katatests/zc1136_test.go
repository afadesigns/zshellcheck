// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1136(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid rm -rf with literal path",
			input:    `rm -rf /tmp/build`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid rm without -rf",
			input:    `rm $file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid rm -rf with variable",
			input: `rm -rf $dir`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1136",
					Message: "Avoid `rm -rf $var` without safeguards. Use `rm -rf ${var:?}` to abort if the variable is empty, preventing accidental deletion.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1136")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
