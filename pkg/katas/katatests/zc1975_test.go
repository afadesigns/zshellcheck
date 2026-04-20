package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1975(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt EXEC` (default on, keeps running)",
			input:    `setopt EXEC`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt NO_EXEC` (restores default)",
			input:    `unsetopt NO_EXEC`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unsetopt EXEC`",
			input: `unsetopt EXEC`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1975",
					Message: "`unsetopt EXEC` stops running commands but keeps parsing — every later line becomes a silent no-op. For syntax checks run `zsh -n script.zsh` from outside the script.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt NO_EXEC`",
			input: `setopt NO_EXEC`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1975",
					Message: "`setopt NO_EXEC` stops running commands but keeps parsing — every later line becomes a silent no-op. For syntax checks run `zsh -n script.zsh` from outside the script.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1975")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
