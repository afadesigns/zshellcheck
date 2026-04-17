package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1403(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — echo $SAVEHIST (Zsh)",
			input:    `echo $SAVEHIST`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $HISTFILESIZE",
			input: `echo $HISTFILESIZE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1403",
					Message: "`$HISTFILESIZE` is Bash-only. Zsh uses `$SAVEHIST` for on-disk history size. Setting `HISTFILESIZE` in Zsh has no effect.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1403")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
