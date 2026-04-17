package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1342(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — find without -empty",
			input:    `find . -type f`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — find -empty",
			input: `find . -empty`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1342",
					Message: "Use Zsh `*(L0)` glob qualifier instead of `find -empty`. Add `.` for regular files only: `*(.L0)`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — find -type f -empty",
			input: `find . -type f -empty -delete`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1342",
					Message: "Use Zsh `*(L0)` glob qualifier instead of `find -empty`. Add `.` for regular files only: `*(.L0)`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1342")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
