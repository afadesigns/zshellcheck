package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestCheckZC1007(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{
			input: `chmod 777 file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1007",
					Message: "Avoid using `chmod 777`. It is a security risk.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			input:    `chmod 755 file.txt`,
			expected: []katas.Violation{},
		},
		{
			input:    `ls -l`,
			expected: []katas.Violation{},
		},
		{
			input: `chmod 777 file1 file2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1007",
					Message: "Avoid using `chmod 777`. It is a security risk.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1007")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
