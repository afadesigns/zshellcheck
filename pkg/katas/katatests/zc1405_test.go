package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1405(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — env -i clean env",
			input:    `env -i cmd`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — env -u VAR",
			input: `env -u DEBUG cmd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1405",
					Message: "Use `(unset VAR; cmd)` subshell instead of `env -u VAR cmd`. In-shell scoping, no external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1405")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
