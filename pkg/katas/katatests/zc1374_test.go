package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1374(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — echo other var",
			input:    `echo $var`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $FUNCNEST expecting depth",
			input: `echo $FUNCNEST`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1374",
					Message: "In Zsh, `$FUNCNEST` is the configured limit, not the current depth. Use `${#funcstack}` for current function nesting depth.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1374")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
