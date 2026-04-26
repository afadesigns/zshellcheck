// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1850(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `ssh -o LogLevel=INFO host`",
			input:    `ssh -o LogLevel=INFO host`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ssh host` (default)",
			input:    `ssh host`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `ssh -o LogLevel=QUIET host`",
			input: `ssh -o LogLevel=QUIET host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1850",
					Message: "`ssh -o LogLevel=QUIET` silences host-key, agent-forward, and canonical-hostname warnings — a MITM event produces no stderr. Keep the default level; capture stderr to a log if you need it clean.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `ssh -oLogLevel=fatal host` (attached)",
			input: `ssh -oLogLevel=fatal host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1850",
					Message: "`ssh -o LogLevel=QUIET` silences host-key, agent-forward, and canonical-hostname warnings — a MITM event produces no stderr. Keep the default level; capture stderr to a log if you need it clean.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1850")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
