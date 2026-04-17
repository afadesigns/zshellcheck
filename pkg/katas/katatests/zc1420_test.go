package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1420(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — chmod 755",
			input:    `chmod 755 file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — chmod +s",
			input: `chmod +s binary`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1420",
					Message: "`chmod +s` / `u+s` / `g+s` sets setuid/setgid — privilege-escalation risk. Prefer sudo policy, capabilities, or containerization.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — chmod 4755",
			input: `chmod 4755 binary`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1420",
					Message: "Numeric mode with leading 4/2/6 sets setuid/setgid — privilege-escalation risk.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1420")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
