package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1104(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "export PATH assignment",
			input: `export PATH=$PATH:/bin`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1104",
					Message: "Use the `path` array instead of manually manipulating the `$PATH` string.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "valid path array assignment",
			input:    `path+=('/usr/local/bin')`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1104")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
