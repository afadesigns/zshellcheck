package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1188(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid export other var",
			input:    `export EDITOR=vim`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid export PATH",
			input: `export PATH=$PATH:/usr/local/bin`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1188",
					Message: "Use `path+=(dir)` instead of `export PATH=$PATH:dir`. Zsh ties the `path` array to `$PATH` for cleaner manipulation.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1188")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
