package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1066(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid loop over command",
			input:    `for i in $(ls); do echo $i; done`,
			expected: []katas.Violation{}, // ZC1066 targets cat only, though ZC1050 might catch ls
		},
		{
			name:     "valid loop over glob",
			input:    `for i in *; do echo $i; done`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid loop over cat",
			input: `for i in $(cat file); do echo $i; done`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1066",
					Message: "Avoid iterating over `cat` output. Use `while read` loop or `($(<file))` for line-based iteration.",
					Line:    1,
					Column:  10,
				},
			},
		},
		{
			name:  "invalid loop over backtick cat",
			input: "for i in `cat file`; do echo $i; done",
			expected: []katas.Violation{
				{
					KataID:  "ZC1066",
					Message: "Avoid iterating over `cat` output. Use `while read` loop or `($(<file))` for line-based iteration.",
					Line:    1,
					Column:  10,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1066")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
