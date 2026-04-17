package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1373(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — env without -0",
			input:    `env VAR=val cmd`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — env -0",
			input: `env -0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1373",
					Message: "Use Zsh `${(0)\"$(<file)\"}` to split NUL-terminated content in-shell. `env -0` is usually followed by `xargs -0` or a read loop — both avoided.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1373")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
