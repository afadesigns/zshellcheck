package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1193(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid rm -f",
			input:    `rm -f file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid rm -i",
			input: `rm -i file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1193",
					Message: "Avoid `rm -i` in scripts — it prompts interactively and will hang in non-interactive execution. Remove `-i` or use explicit checks instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1193")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
