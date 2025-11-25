package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1085(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid quoted array expansion",
			input:    `for i in "${items[@]}"; do echo $i; done`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid quoted variable expansion",
			input:    `for i in "$items"; do echo $i; done`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid glob expansion",
			input:    `for i in *.txt; do echo $i; done`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid command substitution (quoted)",
			input:    `for i in "$(ls)"; do echo $i; done`,
			expected: []katas.Violation{},
		},
		{
			name:     "invalid unquoted variable expansion",
			input:    `for i in $items; do echo $i; done`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1085",
					Message: "Unquoted variable expansion in for loop. This will split on IFS (usually space). Quote it to iterate over lines or array elements.",
					Line:    1,
					Column:  10,
				},
			},
		},
		{
			name:     "invalid unquoted array expansion",
			input:    `for i in ${items[@]}; do echo $i; done`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1085",
					Message: "Unquoted variable expansion in for loop. This will split on IFS (usually space). Quote it to iterate over lines or array elements.",
					Line:    1,
					Column:  10,
				},
			},
		},
		{
			name:     "invalid unquoted command substitution",
			input:    `for i in $(ls); do echo $i; done`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1085",
					Message: "Unquoted variable expansion in for loop. This will split on IFS (usually space). Quote it to iterate over lines or array elements.",
					Line:    1,
					Column:  10,
				},
			},
		},
		{
			name:     "invalid mixed unquoted",
			input:    `for i in start $items end; do echo $i; done`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1085",
					Message: "Unquoted variable expansion in for loop. This will split on IFS (usually space). Quote it to iterate over lines or array elements.",
					Line:    1,
					Column:  16,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1085")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
