package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1842(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt CDABLE_VARS` (explicit default)",
			input:    `unsetopt CDABLE_VARS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NOMATCH` (unrelated)",
			input:    `setopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt CDABLE_VARS`",
			input: `setopt CDABLE_VARS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1842",
					Message: "`setopt CDABLE_VARS` turns a failed `cd NAME` into `cd $NAME` — a typo silently lands in whatever directory the matching variable points to. Keep this in `~/.zshrc`; in scripts use `cd \"$dir\" || exit`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_CDABLE_VARS`",
			input: `unsetopt NO_CDABLE_VARS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1842",
					Message: "`unsetopt NO_CDABLE_VARS` turns a failed `cd NAME` into `cd $NAME` — a typo silently lands in whatever directory the matching variable points to. Keep this in `~/.zshrc`; in scripts use `cd \"$dir\" || exit`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1842")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
