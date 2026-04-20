package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1985(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt SH_FILE_EXPANSION` (default)",
			input:    `unsetopt SH_FILE_EXPANSION`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NO_SH_FILE_EXPANSION`",
			input:    `setopt NO_SH_FILE_EXPANSION`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt SH_FILE_EXPANSION`",
			input: `setopt SH_FILE_EXPANSION`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1985",
					Message: "`setopt SH_FILE_EXPANSION` flips expansion order to POSIX — a `~` or `=cmd` sitting inside a `$VAR` value suddenly resolves, so a user-typed `~other/.cache` escapes into another home. Scope with `emulate -LR sh`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_SH_FILE_EXPANSION`",
			input: `unsetopt NO_SH_FILE_EXPANSION`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1985",
					Message: "`unsetopt NO_SH_FILE_EXPANSION` flips expansion order to POSIX — a `~` or `=cmd` sitting inside a `$VAR` value suddenly resolves, so a user-typed `~other/.cache` escapes into another home. Scope with `emulate -LR sh`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1985")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
