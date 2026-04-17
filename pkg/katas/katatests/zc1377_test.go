package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1377(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — echo $aliases",
			input:    `echo $aliases`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $BASH_ALIASES",
			input: `echo $BASH_ALIASES`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1377",
					Message: "`$BASH_ALIASES` is Bash-only. In Zsh use `$aliases` (assoc array) — same structure, e.g. `print -l ${(kv)aliases}`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1377")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
