package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1925(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt EQUALS` (explicit default)",
			input:    `setopt EQUALS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt NOMATCH` (unrelated)",
			input:    `unsetopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unsetopt EQUALS`",
			input: `unsetopt EQUALS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1925",
					Message: "`unsetopt EQUALS` turns off `=cmd` path expansion and tilde-after-colon — `=python`/`=ls` become literals and `PATH=~/bin:$PATH` stops tilde-expanding. Keep on; scope with `setopt LOCAL_OPTIONS; unsetopt EQUALS` inside a function.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt NO_EQUALS`",
			input: `setopt NO_EQUALS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1925",
					Message: "`setopt NO_EQUALS` turns off `=cmd` path expansion and tilde-after-colon — `=python`/`=ls` become literals and `PATH=~/bin:$PATH` stops tilde-expanding. Keep on; scope with `setopt LOCAL_OPTIONS; unsetopt EQUALS` inside a function.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1925")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
