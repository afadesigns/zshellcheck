// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1533(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "invalid — setsid /usr/local/bin/worker",
			input: `setsid /usr/local/bin/worker`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1533",
					Message: "`setsid` detaches the child from the TTY / session — escapes supervision. Prefer a systemd unit; document a detach if one is genuinely needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — setsid -f /usr/local/bin/worker",
			input: `setsid -f /usr/local/bin/worker`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1533",
					Message: "`setsid` detaches the child from the TTY / session — escapes supervision. Prefer a systemd unit; document a detach if one is genuinely needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1533")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
