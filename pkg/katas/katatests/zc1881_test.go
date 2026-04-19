package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1881(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt MULTIBYTE` (explicit default)",
			input:    `setopt MULTIBYTE`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt NOMATCH` (unrelated)",
			input:    `unsetopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unsetopt MULTIBYTE`",
			input: `unsetopt MULTIBYTE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1881",
					Message: "`unsetopt MULTIBYTE` flips every string op to per-byte math — `${#str}` on an emoji returns 4, substrings slice mid-codepoint, `[[ =~ ]]` Unicode ranges break. Keep the option on; byte-count with `wc -c <<< $str` when truly needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt NO_MULTIBYTE`",
			input: `setopt NO_MULTIBYTE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1881",
					Message: "`setopt NO_MULTIBYTE` flips every string op to per-byte math — `${#str}` on an emoji returns 4, substrings slice mid-codepoint, `[[ =~ ]]` Unicode ranges break. Keep the option on; byte-count with `wc -c <<< $str` when truly needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1881")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
