// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1921(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `systemctl stop myapp`",
			input:    `systemctl stop myapp`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `systemctl kill -s HUP myapp` (graceful reload)",
			input:    `systemctl kill -s HUP myapp`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `systemctl kill -s KILL myapp`",
			input: `systemctl kill -s KILL myapp`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1921",
					Message: "`systemctl kill -s KILL` bypasses `ExecStop=` and `TimeoutStopSec=` — lockfiles, sockets, and shm segments survive and the next restart often fails with \"address already in use\". Use `systemctl stop` or `restart` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `systemctl kill myapp --signal=SIGKILL`",
			input: `systemctl kill myapp --signal=SIGKILL`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1921",
					Message: "`systemctl kill --signal=SIGKILL` bypasses `ExecStop=` and `TimeoutStopSec=` — lockfiles, sockets, and shm segments survive and the next restart often fails with \"address already in use\". Use `systemctl stop` or `restart` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1921")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
