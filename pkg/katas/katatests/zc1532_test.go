// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1532(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — screen -S mysession",
			input:    `screen -S mysession`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — tmux attach",
			input:    `tmux attach-session -t work`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — screen -S name -dm cmd",
			input: `screen -S work -dm /usr/local/bin/worker`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1532",
					Message: "`screen -dm` backgrounds work outside systemd — no journal, no cgroup, common persistence technique. Use a systemd unit instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — tmux new-session -d -s name cmd",
			input: `tmux new-session -d -s work /usr/local/bin/worker`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1532",
					Message: "`tmux new-session -d` backgrounds work outside systemd — no journal, no cgroup, common persistence technique. Use a systemd unit instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1532")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
