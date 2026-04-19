package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1938(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt POSIX_JOBS` (explicit default)",
			input:    `unsetopt POSIX_JOBS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt MONITOR` (unrelated)",
			input:    `setopt MONITOR`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt POSIX_JOBS`",
			input: `setopt POSIX_JOBS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1938",
					Message: "`setopt POSIX_JOBS` scopes `%n` / `fg` / `bg` / `disown` per subshell — parent jobs become invisible inside `(…)`. Leave off; scope POSIX job semantics with `emulate -LR sh` inside a function.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_POSIX_JOBS`",
			input: `unsetopt NO_POSIX_JOBS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1938",
					Message: "`unsetopt NO_POSIX_JOBS` scopes `%n` / `fg` / `bg` / `disown` per subshell — parent jobs become invisible inside `(…)`. Leave off; scope POSIX job semantics with `emulate -LR sh` inside a function.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1938")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
