package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1861(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt OCTAL_ZEROES` (explicit default)",
			input:    `unsetopt OCTAL_ZEROES`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NOMATCH` (unrelated)",
			input:    `setopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt OCTAL_ZEROES`",
			input: `setopt OCTAL_ZEROES`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1861",
					Message: "`setopt OCTAL_ZEROES` reinterprets leading-zero integers as octal — `(( n = 0100 ))` assigns 64 instead of 100, breaking timestamp, phone-prefix, and mode parsing. Keep the option off; use `8#100` when you want explicit octal.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_OCTAL_ZEROES`",
			input: `unsetopt NO_OCTAL_ZEROES`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1861",
					Message: "`unsetopt NO_OCTAL_ZEROES` reinterprets leading-zero integers as octal — `(( n = 0100 ))` assigns 64 instead of 100, breaking timestamp, phone-prefix, and mode parsing. Keep the option off; use `8#100` when you want explicit octal.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1861")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
