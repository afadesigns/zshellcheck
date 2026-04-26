// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1222(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid lsof for files",
			input:    `lsof /var/log/syslog`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid lsof -i",
			input: `lsof -i :8080`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1222",
					Message: "Use `ss -tlnp` instead of `lsof -i` for port checks. `ss` is faster and doesn't require elevated permissions.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1222")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
