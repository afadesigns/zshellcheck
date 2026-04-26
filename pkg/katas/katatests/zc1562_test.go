// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1562(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — env cmd",
			input:    `env cmd`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — env -u TMPDIR cmd",
			input:    `env -u TMPDIR cmd`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — env -u PATH cmd",
			input: `env -u PATH cmd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1562",
					Message: "`env -u PATH` clears a security-relevant variable mid-run. Use `env -i` to sanitise, or set the right value explicitly.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — env -u LD_PRELOAD cmd",
			input: `env -u LD_PRELOAD cmd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1562",
					Message: "`env -u LD_PRELOAD` clears a security-relevant variable mid-run. Use `env -i` to sanitise, or set the right value explicitly.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1562")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
