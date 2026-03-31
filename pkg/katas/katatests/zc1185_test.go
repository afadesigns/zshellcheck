package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1185(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid wc -w with file",
			input:    `wc -w file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid wc -l",
			input:    `wc -l`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid wc -w in pipeline",
			input: `wc -w`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1185",
					Message: "Use Zsh `${#${(z)var}}` for word counting instead of piping through `wc -w`. Parameter expansion avoids spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1185")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
