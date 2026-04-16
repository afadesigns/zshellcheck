package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1341(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — find with other predicate",
			input:    `find . -type f`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — find -executable",
			input: `find . -executable`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1341",
					Message: "Use Zsh `*(.x)` glob qualifier instead of `find -executable`. The `.` restricts to regular files and `x` to the executable bit.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — find -executable with -type f",
			input: `find . -type f -executable`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1341",
					Message: "Use Zsh `*(.x)` glob qualifier instead of `find -executable`. The `.` restricts to regular files and `x` to the executable bit.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1341")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
