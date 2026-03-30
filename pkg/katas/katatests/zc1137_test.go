package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1137(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid mktemp",
			input:    `mktemp /tmp/foo.XXXXXX`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid temp with variable",
			input:    `cat $tmpfile`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid hardcoded tmp path",
			input: `cp data /tmp/myapp.log`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1137",
					Message: "Avoid hardcoded `/tmp/` paths. Use `mktemp` or Zsh `=(cmd)` for temp files to prevent race conditions and symlink attacks.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1137")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
