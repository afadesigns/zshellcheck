package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1255(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid curl -fsSL",
			input:    `curl -fsSL https://example.com/file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid curl without -L",
			input: `curl -s https://example.com/file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1255",
					Message: "Use `curl -L` to follow HTTP redirects. Without `-L`, curl returns redirect responses (301/302) instead of the actual content.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1255")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
