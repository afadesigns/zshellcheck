package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1345(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — find without -perm",
			input:    `find . -type f`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — find -perm 755",
			input: `find . -perm 755`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1345",
					Message: "Use Zsh `*(f:mode:)` glob qualifier instead of `find -perm`. Octal (`*(f:0755:)`) or symbolic (`*(f:u+x:)`) expressions are both supported.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — find -perm -u+x",
			input: `find . -perm -u+x`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1345",
					Message: "Use Zsh `*(f:mode:)` glob qualifier instead of `find -perm`. Octal (`*(f:0755:)`) or symbolic (`*(f:u+x:)`) expressions are both supported.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1345")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
