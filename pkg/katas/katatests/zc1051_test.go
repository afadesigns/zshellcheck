package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1051(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "rm with quoted variable",
			input:    `rm "$file"`,
			expected: []katas.Violation{},
		},
		{
			name:     "rm with literal path",
			input:    `rm /tmp/foo`,
			expected: []katas.Violation{},
		},
		{
			name:     "not rm command",
			input:    `echo $var`,
			expected: []katas.Violation{},
		},
		{
			name:  "rm with unquoted variable",
			input: `rm $file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1051",
					Message: "Unquoted variable in `rm`. Quote it to prevent globbing (e.g. `rm \"$VAR\"`).",
					Line:    1,
					Column:  4,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1051")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
