package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1102(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "sudo redirection",
			input: `sudo echo "foo" > /etc/bar`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1102",
					Message: "Redirection happens before `sudo`. This will likely fail permission checks. Use `| sudo tee`.",
					Line:    1,
					Column:  17,
				},
			},
		},
		{
			name:  "sudo append redirection",
			input: `sudo echo "foo" >> /etc/bar`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1102",
					Message: "Redirection happens before `sudo`. This will likely fail permission checks. Use `| sudo tee`.",
					Line:    1,
					Column:  17,
				},
			},
		},
		{
			name:  "valid sudo usage",
			input: `echo "foo" | sudo tee /etc/bar`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1102")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
