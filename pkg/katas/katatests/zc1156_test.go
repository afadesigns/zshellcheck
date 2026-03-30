package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1156(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid ln -s",
			input:    `ln -s target link`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid ln -sf",
			input:    `ln -sf target link`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid hard link",
			input: `ln target link`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1156",
					Message: "Use `ln -s` for symbolic links instead of hard links. Hard links share inodes and don't work across filesystems.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1156")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
