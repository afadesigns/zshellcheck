package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1043(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "no function definition",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:  "global assignment in function",
			input: "myfunc() { x=1; }",
			expected: []katas.Violation{
				{KataID: "ZC1043", Message: "Variable 'x' is assigned without 'local'. It will be global. Use `local x=1`.", Line: 1, Column: 12},
			},
		},
		{
			name:     "local declaration in function",
			input:    "myfunc() { local x=1; }",
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1043")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
