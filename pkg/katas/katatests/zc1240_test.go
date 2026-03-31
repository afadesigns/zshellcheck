package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1240(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid find with maxdepth and delete",
			input:    `find /tmp -maxdepth 1 -name "*.tmp" -delete`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid find without delete",
			input:    `find . -name "*.log"`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid find -delete without maxdepth",
			input: `find /tmp -name "*.tmp" -delete`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1240",
					Message: "Use `find -maxdepth N` with `-delete` to limit deletion scope. Without depth limits, find recurses infinitely.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1240")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
