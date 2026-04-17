package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1409(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — test -f file",
			input:    `test -f file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — test -N file",
			input: `test -N file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1409",
					Message: "`test -N file` (modified-since-read) is a Bash extension. In Zsh use `zmodload zsh/stat; zstat -H s file; (( s[mtime] > s[atime] ))`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1409")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
