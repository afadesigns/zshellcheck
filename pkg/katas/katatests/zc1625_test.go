package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1625(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — rm -rf scoped path",
			input:    `rm -rf /tmp/staging`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — rm with --preserve-root=all",
			input:    `rm -rf --preserve-root=all $TARGET`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — rm -rf --no-preserve-root /",
			input: `rm -rf --no-preserve-root /`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1625",
					Message: "`rm --no-preserve-root` disables the GNU safeguard against `rm -rf /`. Remove the flag; if a specific path needs deletion, list it explicitly.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — rm -rf --no-preserve-root $TARGET",
			input: `rm -rf --no-preserve-root $TARGET`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1625",
					Message: "`rm --no-preserve-root` disables the GNU safeguard against `rm -rf /`. Remove the flag; if a specific path needs deletion, list it explicitly.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1625")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
