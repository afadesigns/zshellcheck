package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1942(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt CLOBBER_EMPTY` (explicit default)",
			input:    `unsetopt CLOBBER_EMPTY`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NO_CLOBBER` (unrelated)",
			input:    `setopt NO_CLOBBER`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt CLOBBER_EMPTY`",
			input: `setopt CLOBBER_EMPTY`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1942",
					Message: "`setopt CLOBBER_EMPTY` lets `>file` overwrite zero-length files even under `NO_CLOBBER` — `touch`ed lock / sentinel files lose their safety net. Keep off; use explicit `>|file` to bypass `NO_CLOBBER` for a specific write.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_CLOBBER_EMPTY`",
			input: `unsetopt NO_CLOBBER_EMPTY`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1942",
					Message: "`unsetopt NO_CLOBBER_EMPTY` lets `>file` overwrite zero-length files even under `NO_CLOBBER` — `touch`ed lock / sentinel files lose their safety net. Keep off; use explicit `>|file` to bypass `NO_CLOBBER` for a specific write.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1942")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
