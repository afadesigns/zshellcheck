package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1397(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — echo $compstate (Zsh)",
			input:    `echo $compstate`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $COMP_TYPE",
			input: `echo $COMP_TYPE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1397",
					Message: "Bash `$COMP_TYPE`/`$COMP_KEY`/`$COMP_WORDBREAKS` are not Zsh-native. Use `$compstate` associative array for completion context in Zsh.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1397")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
