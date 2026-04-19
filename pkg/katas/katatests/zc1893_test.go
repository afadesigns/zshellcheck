package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1893(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt BARE_GLOB_QUAL` (explicit default)",
			input:    `setopt BARE_GLOB_QUAL`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt NOMATCH` (unrelated)",
			input:    `unsetopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unsetopt BARE_GLOB_QUAL`",
			input: `unsetopt BARE_GLOB_QUAL`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1893",
					Message: "`unsetopt BARE_GLOB_QUAL` disables `*(qualifier)` syntax — `*(N)` stops being null-glob and becomes an alternation, so null-glob idioms silently break. Keep the option on; scope inside a `LOCAL_OPTIONS` function.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt NO_BARE_GLOB_QUAL`",
			input: `setopt NO_BARE_GLOB_QUAL`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1893",
					Message: "`setopt NO_BARE_GLOB_QUAL` disables `*(qualifier)` syntax — `*(N)` stops being null-glob and becomes an alternation, so null-glob idioms silently break. Keep the option on; scope inside a `LOCAL_OPTIONS` function.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1893")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
