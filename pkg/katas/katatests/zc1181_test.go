package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1181(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid open with -a flag",
			input:    `open -a Safari`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid xdg-open",
			input: `xdg-open https://example.com`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1181",
					Message: "Use `$BROWSER` or check `$OSTYPE` instead of `xdg-open` for portable URL/file opening across Linux and macOS.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid open URL",
			input: `open https://example.com`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1181",
					Message: "Use `$BROWSER` or check `$OSTYPE` instead of `open` for portable URL/file opening across Linux and macOS.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1181")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
