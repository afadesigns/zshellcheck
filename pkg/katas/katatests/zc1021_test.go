package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestCheckZC1021(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{
			input: `chmod 755 file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1021",
					Message: "Use symbolic permissions with `chmod` instead of octal.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			input:    `chmod u+x file`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1021")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
