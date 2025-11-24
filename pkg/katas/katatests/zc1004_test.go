package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1004(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid return",
			input:    `my_func() { return 0; }`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid exit in subshell",
			input:    `my_func() { ( exit 1 ) }`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid exit in command sub",
			input:    `my_func() { local x=$(exit 1); }`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid exit",
			input: `my_func() { exit 1; }`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1004",
					Message: "Use `return` instead of `exit` in functions to avoid killing the shell.",
					Line:    1,
					Column:  13,
				},
			},
		},
		{
			name:  "invalid exit deep",
			input: `my_func() { if true; then exit 1; fi }`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1004",
					Message: "Use `return` instead of `exit` in functions to avoid killing the shell.",
					Line:    1,
					Column:  27,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1004")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
