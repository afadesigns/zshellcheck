package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1074(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid modifier usage",
			input:    `echo ${path:h} ${file:t}`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid dirname usage",
			input: `dir=$(dirname $path)`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1074",
					Message: "Use '${var:h}' instead of '$(dirname $var)'. Modifiers are faster and built-in.",
					Line:    1,
					Column:  5,
				},
			},
		},
		{
			name:  "invalid basename usage",
			input: `base=$(basename $path)`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1074",
					Message: "Use '${var:t}' instead of '$(basename $var)'. Modifiers are faster and built-in.",
					Line:    1,
					Column:  6,
				},
			},
		},
		{
			name:  "invalid backtick dirname",
			input: "dir=`dirname $path`",
			expected: []katas.Violation{
				{
					KataID:  "ZC1074",
					Message: "Use '${var:h}' instead of '$(dirname $var)'. Modifiers are faster and built-in.",
					Line:    1,
					Column:  5,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1074")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
