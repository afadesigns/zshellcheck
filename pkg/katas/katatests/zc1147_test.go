package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1147(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid mkdir -p",
			input:    `mkdir -p /tmp/a/b/c`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid mkdir single level",
			input:    `mkdir newdir`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid mkdir nested without -p",
			input: `mkdir /tmp/a/b/c`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1147",
					Message: "Use `mkdir -p` when creating nested directories. Without `-p`, `mkdir` fails if parent directories don't exist.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1147")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
