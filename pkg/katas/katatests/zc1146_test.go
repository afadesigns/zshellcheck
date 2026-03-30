package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1146(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid awk with file",
			input:    `awk '{print}' file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid cat with flags piped",
			input:    `cat -n file | awk '{print}'`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid cat | awk",
			input: `cat data.csv | awk -F, '{print $1}'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1146",
					Message: "Pass the file directly to `awk` instead of `cat file | awk`. Most text-processing tools accept file arguments.",
					Line:    1,
					Column:  14,
				},
			},
		},
		{
			name:  "invalid cat | sort",
			input: `cat names.txt | sort`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1146",
					Message: "Pass the file directly to `sort` instead of `cat file | sort`. Most text-processing tools accept file arguments.",
					Line:    1,
					Column:  15,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1146")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
