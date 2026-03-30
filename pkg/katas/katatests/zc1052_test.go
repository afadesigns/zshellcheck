package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1052(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "sed without -i",
			input:    `sed 's/foo/bar/' file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "not sed command",
			input:    `grep foo bar`,
			expected: []katas.Violation{},
		},
		{
			name:     "sed with other flags",
			input:    `sed -e 's/foo/bar/' file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "sed with -i",
			input: `sed -i 's/foo/bar/' file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1052",
					Message: "`sed -i` is non-portable (GNU vs BSD differences). Use `perl -i` or a temporary file.",
					Line:    1,
					Column:  5,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1052")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
