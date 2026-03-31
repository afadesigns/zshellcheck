package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1241(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid xargs -0 rm",
			input:    `xargs -0 rm`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid xargs without rm",
			input:    `xargs grep pattern`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid xargs rm without -0",
			input: `xargs rm`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1241",
					Message: "Use `xargs -0 rm` with `find -print0` for safe deletion. Without `-0`, filenames with spaces or special characters break.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1241")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
