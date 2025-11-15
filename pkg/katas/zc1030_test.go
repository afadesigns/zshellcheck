package katas

import (
	"testing"
)

func TestCheckZC1030(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Violation
	}{
		{
			name:  "echo with a simple string",
			input: `echo "hello"`,
			expected: []Violation{
				{
					KataID:  "ZC1030",
					Message: "Use `printf` for more reliable and portable string formatting instead of `echo`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "printf with a simple string",
			input:    `printf "hello"`,
			expected: []Violation{},
		},
		{
			name:     "echo with a variable",
			input:    `echo "$foo"`,
			expected: []Violation{
				{
					KataID:  "ZC1030",
					Message: "Use `printf` for more reliable and portable string formatting instead of `echo`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := check(tt.input, "ZC1030")
			assertViolations(t, tt.input, violations, tt.expected)
		})
	}
}