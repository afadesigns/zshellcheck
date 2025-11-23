package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1072(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid awk",
			input:    `awk '/pattern/ {print}' file`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid grep recursive",
			input:    `grep -r pattern . | awk '{print $1}'`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid grep awk",
			input: `grep pattern file | awk '{print $1}'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1072",
					Message: "Use `awk '/pattern/ {...}'` instead of `grep pattern | awk '{...}'` to avoid a pipeline.",
					Line:    1,
					Column:  19, // Position of pipe
				},
			},
		},
		{
			name:  "invalid grep awk with flags",
			input: `grep -i pattern file | awk '{print}'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1072",
					Message: "Use `awk '/pattern/ {...}'` instead of `grep pattern | awk '{...}'` to avoid a pipeline.",
					Line:    1,
					Column:  22,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1072")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
