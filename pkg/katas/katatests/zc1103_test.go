package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1103(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "PATH assignment",
			input: `PATH=$PATH:/usr/local/bin`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1103",
					Message: "Use the `path` array instead of manually manipulating the `$PATH` string.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "export PATH",
			input: `export PATH=$PATH:/bin`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1103",
					Message: "Use the `path` array instead of manually manipulating the `$PATH` string.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "path array assignment",
			input: `path+=('/usr/local/bin')`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1103")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
