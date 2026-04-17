package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1372(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — zmv usage",
			input:    `zmv '(*).txt' '$1.md'`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — rename perl-style",
			input: `rename 's/\.txt$/.md/' *.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1372",
					Message: "Use Zsh `zmv` (autoload -Uz zmv) instead of `rename`/`rename.ul`/`prename`. Glob-pattern renaming is handled in-shell with capture groups.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — rename.ul util-linux",
			input: `rename.ul .txt .md *.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1372",
					Message: "Use Zsh `zmv` (autoload -Uz zmv) instead of `rename`/`rename.ul`/`prename`. Glob-pattern renaming is handled in-shell with capture groups.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1372")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
