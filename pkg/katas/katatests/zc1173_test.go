package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1173(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid column -s",
			input:    `column -s: file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid column -t",
			input: `column -t file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1173",
					Message: "Use Zsh `print -C N` for columnar output instead of `column -t`. The `print` builtin formats columns without spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1173")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
