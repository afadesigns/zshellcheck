package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1048(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "source with absolute path",
			input:    `source /etc/profile`,
			expected: []katas.Violation{},
		},
		{
			name:     "not source command",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:     "source with no arguments",
			input:    `source`,
			expected: []katas.Violation{},
		},
		{
			name:  "source with relative path",
			input: `source ./lib.zsh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1048",
					Message: "Avoid `source` with relative paths. Use `${0:a:h}/...` to resolve relative to the script.",
					Line:    1,
					Column:  8,
				},
			},
		},
		{
			name:  "dot with relative path",
			input: `. ../lib.zsh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1048",
					Message: "Avoid `source` with relative paths. Use `${0:a:h}/...` to resolve relative to the script.",
					Line:    1,
					Column:  3,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1048")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
