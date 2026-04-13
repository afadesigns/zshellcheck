package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1323(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid kill usage",
			input:    `kill -STOP 1234`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid suspend usage",
			input: `suspend -f`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1323",
					Message: "Avoid `suspend` in Zsh — it is a Bash builtin. Use `kill -STOP $$` if needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1323")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
