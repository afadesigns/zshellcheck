package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1198(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid sed -i",
			input:    `sed -i 's/old/new/' file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid nano in script",
			input: `nano config.yaml`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1198",
					Message: "Avoid `nano` in scripts — interactive editors hang without a terminal. Use `sed -i` or `ed` for scripted file editing.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid vim in script",
			input: `vim /etc/hosts`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1198",
					Message: "Avoid `vim` in scripts — interactive editors hang without a terminal. Use `sed -i` or `ed` for scripted file editing.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1198")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
