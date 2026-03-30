package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1055(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "no comparison",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:  "compare with empty double-quoted string",
			input: `[[ $x == "" ]]`,
			expected: []katas.Violation{
				{KataID: "ZC1055", Message: "Use `[[ -z ... ]]` instead of comparing with empty string.", Line: 1, Column: 8},
			},
		},
		{
			name:  "compare not equal empty string",
			input: `[[ $x != "" ]]`,
			expected: []katas.Violation{
				{KataID: "ZC1055", Message: "Use `[[ -n ... ]]` instead of comparing with empty string.", Line: 1, Column: 8},
			},
		},
		{
			name:     "compare with non-empty string",
			input:    `[[ $x == "hello" ]]`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1055")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
