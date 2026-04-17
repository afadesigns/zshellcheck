package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1353(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — printf without -v",
			input:    `printf 'hello %s\n' world`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — printf -v var",
			input: `printf -v line '%d' 42`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1353",
					Message: "Avoid `printf -v` in Zsh — use `print -v var -rf fmt ...` or `var=$(printf fmt ...)`. `-v` is Bash-specific and ignored elsewhere.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1353")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
