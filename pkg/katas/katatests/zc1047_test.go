package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1047(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "sudo command",
			input: `sudo apt install vim`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1047",
					Message: "Avoid `sudo` in scripts. Run the entire script as root if privileges are required.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "not sudo",
			input:    `apt install vim`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1047")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
