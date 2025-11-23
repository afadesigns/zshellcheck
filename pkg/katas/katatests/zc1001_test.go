package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1001(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid array access",
			input:    `echo ${my_array[1]}`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid array access",
			input: `echo $my_array[1]`,
			expected: []katas.Violation{
				{
					KataID: "ZC1001",
					Message: "Use ${} for array element access. " +
						"Accessing array elements with `$my_array[...]` is not the correct syntax in Zsh.",
					Line:   1,
					Column: 6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1001")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
