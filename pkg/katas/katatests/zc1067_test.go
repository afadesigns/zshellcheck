package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1067(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid export separation",
			input:    `var=$(cmd); export var`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid export literal",
			input:    `export var="value"`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid export command substitution",
			input: `export var=$(cmd)`,
			expected: []katas.Violation{
				{
					KataID: "ZC1067",
					Message: "Exporting and assigning a command substitution in one step masks the return value. " +
						"Use `var=$(cmd); export var`.",
					Line:   1,
					Column: 1,
				},
			},
		},
		{
			name:  "invalid export backticks",
			input: "export var=`cmd`",
			expected: []katas.Violation{
				{
					KataID: "ZC1067",
					Message: "Exporting and assigning a command substitution in one step masks the return value. " +
						"Use `var=$(cmd); export var`.",
					Line:   1,
					Column: 1,
				},
			},
		},
		{
			name:     "valid export with no assignment",
			input:    `export var`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1067")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
