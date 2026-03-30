package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1117(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "invalid nohup usage",
			input: `nohup ./server`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1117",
					Message: "Use `cmd &!` or `cmd & disown` instead of `nohup cmd &`. Zsh `&!` is a built-in shorthand that avoids spawning nohup.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid nohup with redirect",
			input: `nohup ./server > /dev/null`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1117",
					Message: "Use `cmd &!` or `cmd & disown` instead of `nohup cmd &`. Zsh `&!` is a built-in shorthand that avoids spawning nohup.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1117")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
