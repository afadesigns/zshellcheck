package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1391(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — test -n VAR",
			input:    `test -n "$VAR"`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — test -v VAR",
			input: `test -v VAR`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1391",
					Message: "Use `(( ${+VAR} ))` for Zsh set-check — `-v` is a Bash 4.2+ extension, not reliably portable to Zsh.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1391")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
