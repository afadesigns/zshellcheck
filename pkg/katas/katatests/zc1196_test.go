package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1196(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid less with flag",
			input:    `less -R file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid less in script",
			input: `less file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1196",
					Message: "Avoid `less` in scripts — it requires interactive terminal input. Use `cat` or redirect output to a pager only when `$TERM` is available.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1196")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
