package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1845(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt PATH_DIRS` (explicit default)",
			input:    `unsetopt PATH_DIRS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NOMATCH` (unrelated)",
			input:    `setopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt PATH_DIRS`",
			input: `setopt PATH_DIRS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1845",
					Message: "`setopt PATH_DIRS` lets `subdir/cmd` fall back to a `$PATH` lookup — a missing local binary silently runs a same-named subtree elsewhere on `$PATH`. Leave the option off; call locals as `./subdir/cmd`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_PATH_DIRS`",
			input: `unsetopt NO_PATH_DIRS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1845",
					Message: "`unsetopt NO_PATH_DIRS` lets `subdir/cmd` fall back to a `$PATH` lookup — a missing local binary silently runs a same-named subtree elsewhere on `$PATH`. Leave the option off; call locals as `./subdir/cmd`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1845")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
