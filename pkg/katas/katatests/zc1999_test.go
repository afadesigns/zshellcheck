package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1999(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt AUTO_NAMED_DIRS` (default)",
			input:    `unsetopt AUTO_NAMED_DIRS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NO_AUTO_NAMED_DIRS`",
			input:    `setopt NO_AUTO_NAMED_DIRS`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt AUTO_NAMED_DIRS`",
			input: `setopt AUTO_NAMED_DIRS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1999",
					Message: "`setopt AUTO_NAMED_DIRS` auto-registers every dir-valued scalar as `~name` — collisions with real usernames and stray `~$var` expansions. Register named dirs explicitly with `hash -d NAME=PATH`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_AUTO_NAMED_DIRS`",
			input: `unsetopt NO_AUTO_NAMED_DIRS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1999",
					Message: "`unsetopt NO_AUTO_NAMED_DIRS` auto-registers every dir-valued scalar as `~name` — collisions with real usernames and stray `~$var` expansions. Register named dirs explicitly with `hash -d NAME=PATH`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1999")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
