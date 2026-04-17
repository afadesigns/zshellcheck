package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1384(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — unrelated echo",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $EXECIGNORE",
			input: `echo $EXECIGNORE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1384",
					Message: "`$EXECIGNORE` is Bash-only. For completion filtering in Zsh use `zstyle ':completion:*' ignored-patterns 'pattern'`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1384")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
