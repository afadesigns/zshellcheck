package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1937(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `tmux list-sessions`",
			input:    `tmux list-sessions`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `tmux kill-window -t dev:1`",
			input:    `tmux kill-window -t dev:1`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `tmux kill-server`",
			input: `tmux kill-server`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1937",
					Message: "`tmux kill-server` tears down every detached process inside the session — builds, log tails, port-forwards get `SIGHUP`'d with no cleanup. Use `kill-window` for surgical removal or `systemd-run --scope` for workloads.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `screen -X quit`",
			input: `screen -X quit`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1937",
					Message: "`screen -X quit` tears down every detached process inside the session — builds, log tails, port-forwards get `SIGHUP`'d with no cleanup. Use `kill-window` for surgical removal or `systemd-run --scope` for workloads.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1937")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
