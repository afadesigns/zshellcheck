package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestCheckZC1038(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{
			input: `cat file | grep "foo"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1038",
					Message: "Avoid useless use of cat. Prefer `command file` or `command < file` over `cat file | command`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			input:    `cat file1 file2 | grep "foo"`, // Valid concatenation
			expected: []katas.Violation{},
		},
		{
			input:    `grep "foo" file`, // Direct file access
			expected: []katas.Violation{},
		},
		{
			input:    `cat | grep "foo"`, // Reading from stdin
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1038")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
