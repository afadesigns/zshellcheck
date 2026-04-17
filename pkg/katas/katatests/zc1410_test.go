package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1410(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — zsh plain",
			input:    `zsh script.zsh`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — bash --rcfile",
			input: `bash --rcfile myrc script.sh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1410",
					Message: "`bash --rcfile`/`--init-file` have no Zsh equivalent flag. Use `ZDOTDIR=/path zsh` to relocate all Zsh rc files.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1410")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
