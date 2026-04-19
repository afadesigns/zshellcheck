package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1908(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt MAGIC_EQUAL_SUBST` (explicit default)",
			input:    `unsetopt MAGIC_EQUAL_SUBST`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt EXTENDED_GLOB` (unrelated)",
			input:    `setopt EXTENDED_GLOB`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt MAGIC_EQUAL_SUBST`",
			input: `setopt MAGIC_EQUAL_SUBST`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1908",
					Message: "`setopt MAGIC_EQUAL_SUBST` gives every `key=value` argument tilde/parameter expansion on the RHS — literal CLI args like `rsync host:dst=~/backup` silently change. Keep it off; quote the assignment if expansion is really wanted.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_MAGIC_EQUAL_SUBST`",
			input: `unsetopt NO_MAGIC_EQUAL_SUBST`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1908",
					Message: "`unsetopt NO_MAGIC_EQUAL_SUBST` gives every `key=value` argument tilde/parameter expansion on the RHS — literal CLI args like `rsync host:dst=~/backup` silently change. Keep it off; quote the assignment if expansion is really wanted.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1908")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
