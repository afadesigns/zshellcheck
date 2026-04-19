package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1920(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt VERBOSE` (explicit default)",
			input:    `unsetopt VERBOSE`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt EXTENDED_HISTORY` (unrelated)",
			input:    `setopt EXTENDED_HISTORY`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt VERBOSE`",
			input: `setopt VERBOSE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1920",
					Message: "`setopt VERBOSE` echoes every executed command to stderr — any line that mentions a password, token, or API key leaks with the trace. Remove and use `printf` / a logger, or scope via `setopt LOCAL_OPTIONS VERBOSE` in a helper.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_VERBOSE`",
			input: `unsetopt NO_VERBOSE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1920",
					Message: "`unsetopt NO_VERBOSE` echoes every executed command to stderr — any line that mentions a password, token, or API key leaks with the trace. Remove and use `printf` / a logger, or scope via `setopt LOCAL_OPTIONS VERBOSE` in a helper.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1920")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
