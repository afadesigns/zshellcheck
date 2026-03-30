package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1113(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid readlink without resolve flag",
			input:    `readlink /path/to/link`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid realpath",
			input: `realpath /path/to/file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1113",
					Message: "Use `${var:A}` instead of `realpath` to resolve absolute paths. Zsh path modifiers avoid spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid readlink -f",
			input: `readlink -f /path/to/link`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1113",
					Message: "Use `${var:A}` instead of `readlink -f` to resolve absolute paths. Zsh path modifiers avoid spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "valid realpath with complex flags",
			input:    `realpath --relative-to /base /path`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1113")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
