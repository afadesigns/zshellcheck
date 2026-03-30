package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1115(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid rev with file",
			input:    `rev file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid rev in pipeline",
			input: `rev -`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1115",
					Message: "Use Zsh string manipulation instead of `rev`. Parameter expansion can reverse strings without spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1115")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
