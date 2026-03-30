package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1178(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid stty raw",
			input:    `stty raw`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid stty size",
			input: `stty size`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1178",
					Message: "Use `$COLUMNS` and `$LINES` instead of `stty size`. Zsh tracks terminal dimensions as built-in variables.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1178")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
