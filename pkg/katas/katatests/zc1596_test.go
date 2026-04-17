package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1596(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — emulate -L sh",
			input:    `emulate -L sh`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — emulate -LR bash",
			input:    `emulate -LR bash`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — emulate zsh (reset to zsh)",
			input:    `emulate zsh`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — emulate sh",
			input: `emulate sh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1596",
					Message: "`emulate sh` without `-L` flips the options for the whole shell. Use `emulate -L sh` inside a function, or rename the script to `.sh` if Zsh features are not needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — emulate -R bash",
			input: `emulate -R bash`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1596",
					Message: "`emulate bash` without `-L` flips the options for the whole shell. Use `emulate -L bash` inside a function, or rename the script to `.sh` if Zsh features are not needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1596")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
