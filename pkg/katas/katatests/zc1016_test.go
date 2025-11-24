package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1016(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "safe read",
			input:    `read name`,
			expected: []katas.Violation{},
		},
		{
			name:     "safe read password with -s",
			input:    `read -s password`,
			expected: []katas.Violation{},
		},
		{
			name:     "safe read with combined flags",
			input:    `read -rs password`,
			expected: []katas.Violation{},
		},
		{
			name:  "unsafe read password",
			input: `read password`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1016",
					Message: "Use `read -s` to hide input when reading sensitive variable 'password'.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "unsafe read with prompt",
			input: `read "secret_key?Enter key: "`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1016",
					Message: "Use `read -s` to hide input when reading sensitive variable 'secret_key'.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "unsafe read multiple vars",
			input: `read user password`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1016",
					Message: "Use `read -s` to hide input when reading sensitive variable 'password'.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1016")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
