package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1376(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — unrelated echo",
			input:    `echo $VAR`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $BASH_XTRACEFD",
			input: `echo $BASH_XTRACEFD`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1376",
					Message: "`BASH_XTRACEFD` is Bash-only. Zsh ignores it. Redirect trace output with `exec {fd}>file; exec 2>&$fd; setopt XTRACE` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1376")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
