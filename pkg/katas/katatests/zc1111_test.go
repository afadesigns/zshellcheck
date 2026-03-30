package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1111(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid xargs with null separator",
			input:    `xargs -0 rm`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid xargs with parallel",
			input:    `xargs -P 4 cmd`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid xargs with replace",
			input:    `xargs -I {} mv {} /dest`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid simple xargs",
			input: `xargs rm`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1111",
					Message: "Consider using Zsh array iteration instead of `xargs`. `for item in ${(f)$(cmd)}` splits output by newlines without spawning xargs.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid xargs with command only",
			input: `xargs grep pattern`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1111",
					Message: "Consider using Zsh array iteration instead of `xargs`. `for item in ${(f)$(cmd)}` splits output by newlines without spawning xargs.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1111")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
