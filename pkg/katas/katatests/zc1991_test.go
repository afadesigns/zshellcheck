package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1991(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt CSH_NULLCMD` (default)",
			input:    `unsetopt CSH_NULLCMD`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NO_CSH_NULLCMD`",
			input:    `setopt NO_CSH_NULLCMD`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt CSH_NULLCMD`",
			input: `setopt CSH_NULLCMD`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1991",
					Message: "`setopt CSH_NULLCMD` makes `> file` / `< file` (no command) a parse error — log truncation and bare-redirect idioms stop working. Write `: > file` explicitly for truncation.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_CSH_NULLCMD`",
			input: `unsetopt NO_CSH_NULLCMD`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1991",
					Message: "`unsetopt NO_CSH_NULLCMD` makes `> file` / `< file` (no command) a parse error — log truncation and bare-redirect idioms stop working. Write `: > file` explicitly for truncation.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1991")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
