package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1347(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — find without -group",
			input:    `find . -type f`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — find -group",
			input: `find . -group wheel`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1347",
					Message: "Use Zsh `*(g:name:)` / `*(g+gid)` / `*(G)` glob qualifiers instead of `find -group`/`-gid`/`-nogroup`. Group predicates live entirely in the shell.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — find -gid 10",
			input: `find . -gid 10`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1347",
					Message: "Use Zsh `*(g:name:)` / `*(g+gid)` / `*(G)` glob qualifiers instead of `find -group`/`-gid`/`-nogroup`. Group predicates live entirely in the shell.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1347")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
