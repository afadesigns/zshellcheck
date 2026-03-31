package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1194(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid single sed expression",
			input:    `sed -e 's/foo/bar/' file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid multiple sed -e",
			input: `sed -e 's/foo/bar/' -e 's/baz/qux/' file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1194",
					Message: "Combine multiple `sed -e` expressions into a single script: `sed 's/a/b/; s/c/d/'` is cleaner than multiple `-e` flags.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1194")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
