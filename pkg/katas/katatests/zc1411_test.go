package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1411(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — enable builtin",
			input:    `enable name`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — enable -n builtin",
			input: `enable -n echo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1411",
					Message: "Use Zsh `disable name` instead of `enable -n name`. Zsh has a dedicated `disable` builtin that reads clearer.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1411")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
