package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1404(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — echo $commands (Zsh)",
			input:    `echo $commands`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $BASH_CMDS",
			input: `echo $BASH_CMDS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1404",
					Message: "`$BASH_CMDS` is Bash-only. In Zsh use `$commands` (assoc array, names→paths) via `zsh/parameter`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1404")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
