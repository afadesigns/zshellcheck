package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1132(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid grep -o with file",
			input:    `grep -o pattern file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid grep without -o",
			input:    `grep pattern`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid grep -o in pipeline",
			input: `grep -o pattern`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1132",
					Message: "Use Zsh pattern extraction `${(M)var:#pattern}` or `[[ $var =~ regex ]] && echo $match[1]` instead of piping through `grep -o`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1132")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
