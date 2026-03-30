package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1093(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid redirect",
			input:    `grep pattern < file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid cat with flags",
			input:    `cat -n file.txt | head`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid cat multiple files",
			input:    `cat file1 file2 | sort`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid useless cat",
			input: `cat file.txt | grep pattern`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1093",
					Message: "`cat file | command` is inefficient. Use `command < file` or pass the filename as an argument.",
					Line:    1,
					Column:  14,
				},
			},
		},
		{
			name:  "invalid useless cat with sort",
			input: `cat data.csv | sort`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1093",
					Message: "`cat file | command` is inefficient. Use `command < file` or pass the filename as an argument.",
					Line:    1,
					Column:  14,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1093")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
