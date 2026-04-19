package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1865(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt CASE_MATCH` (explicit default)",
			input:    `setopt CASE_MATCH`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt NOMATCH` (unrelated)",
			input:    `unsetopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unsetopt CASE_MATCH`",
			input: `unsetopt CASE_MATCH`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1865",
					Message: "`unsetopt CASE_MATCH` flips every `[[ =~ ]]` / `[[ == pat ]]` to case-insensitive — `Admin` matches `ADMIN`, dispatchers collide. Keep it on; scope per-line with `(#i)pattern`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt NO_CASE_MATCH`",
			input: `setopt NO_CASE_MATCH`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1865",
					Message: "`setopt NO_CASE_MATCH` flips every `[[ =~ ]]` / `[[ == pat ]]` to case-insensitive — `Admin` matches `ADMIN`, dispatchers collide. Keep it on; scope per-line with `(#i)pattern`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1865")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
