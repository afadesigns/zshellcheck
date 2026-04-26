// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1426(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — git clone https URL",
			input:    `git clone https://github.com/owner/repo.git`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — git clone http URL",
			input: `git clone http://example.com/repo.git`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1426",
					Message: "`git clone http://` is unencrypted/unauthenticated. Use `https://` or SSH with verified host keys.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1426")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
