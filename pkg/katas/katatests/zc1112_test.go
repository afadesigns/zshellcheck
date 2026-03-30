package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1112(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid grep -c with file",
			input:    `grep -c pattern file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid grep without -c",
			input:    `grep pattern`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid grep -c in pipeline",
			input: `grep -c pattern`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1112",
					Message: "Use Zsh array filtering `${(M)array:#pattern}` or `${#${(f)...}}` for counting instead of `grep -c`. Avoids spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "valid grep -v (not count)",
			input:    `grep -v pattern`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1112")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
