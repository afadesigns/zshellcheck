package katas

import (
	"testing"
)

func TestCheckZC1037(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Violation
	}{
		{
			name:     "echo with a simple string should not trigger",
			input:    `echo "hello"`,
			expected: []Violation{},
		},
		{
			name:  "echo with an unquoted variable should trigger",
			input: `echo $foo`,
			expected: []Violation{
				{
					KataID:  "ZC1037",
					Message: "Use 'print -r --' instead of 'echo' to reliably print variable expansions.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := check(tt.input, "ZC1037")
			assertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
