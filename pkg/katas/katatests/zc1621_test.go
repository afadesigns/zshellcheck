package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1621(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — default tmux",
			input:    `tmux new-session -d`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — tmux -S in XDG_RUNTIME_DIR",
			input:    `tmux -S $XDG_RUNTIME_DIR/tmux-main new-session`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — tmux -S /tmp/shared",
			input: `tmux -S /tmp/shared new-session`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1621",
					Message: "`tmux -S /tmp/shared` places the socket in a world-traversable directory — any local user who can read the socket can attach the session. Use `$XDG_RUNTIME_DIR` or a 0700-scoped parent dir.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — tmux -S /var/tmp/pair",
			input: `tmux -S /var/tmp/pair attach`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1621",
					Message: "`tmux -S /var/tmp/pair` places the socket in a world-traversable directory — any local user who can read the socket can attach the session. Use `$XDG_RUNTIME_DIR` or a 0700-scoped parent dir.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1621")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
