package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1087(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid redirection to new file",
			input:    `cat input.txt > output.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid append redirection",
			input:    `cat file.txt >> file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "invalid clobbering simple command",
			input:    `sort file.txt > file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1087",
					Message: "Output redirection overwrites input file `file.txt`. The file is truncated before reading.",
					Line:    1,
					Column:  15, // Points to >
				},
			},
		},
		{
			name:     "invalid clobbering pipeline",
			input:    `cat file.txt | grep foo > file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1087",
					Message: "Output redirection overwrites input file `file.txt`. The file is truncated before reading.",
					Line:    1,
					Column:  25, // Points to >
				},
			},
		},
		{
			name:     "invalid clobbering with input redirection",
			input:    `grep foo < file.txt > file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1087",
					Message: "Output redirection overwrites input file `file.txt`. The file is truncated before reading.",
					Line:    1,
					Column:  21, // Points to >
				},
			},
		},
		{
			name:     "valid different files",
			input:    `sed 's/a/b/' input > output`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1087")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
