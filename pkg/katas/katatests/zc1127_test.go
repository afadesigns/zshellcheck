package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1127(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid ls -la",
			input:    `ls -la`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid ls -1",
			input: `ls -1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1127",
					Message: "Use Zsh glob qualifiers `files=(*(N)); echo ${#files}` instead of `ls -1 | wc -l`. Avoids spawning external processes for file counting.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1127")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
