package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1371(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — basename single path",
			input:    `basename /usr/bin/zsh`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — basename -a",
			input: `basename -a /a/b /c/d`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1371",
					Message: "Use Zsh `${paths:t}` on an array for bulk basename extraction instead of `basename -a`. The `:t` modifier applies to every array element.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1371")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
