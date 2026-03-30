package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1135(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid env -i",
			input:    `env -i cmd`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid env VAR=val cmd",
			input: `env FOO=bar cmd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1135",
					Message: "Use inline `VAR=val cmd` instead of `env VAR=val cmd`. Zsh supports inline env assignment without spawning env.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1135")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
