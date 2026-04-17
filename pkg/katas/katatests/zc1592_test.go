package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1592(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — faillock status",
			input:    `faillock -u bob`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — pam_tally2 status",
			input:    `pam_tally2 -u bob`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — faillock -u bob --reset",
			input: `faillock -u bob --reset`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1592",
					Message: "`faillock --reset` clears the PAM failed-auth counter — masks ongoing brute force. Log the prior count and alert before resetting.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — pam_tally2 -r -u bob",
			input: `pam_tally2 -r -u bob`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1592",
					Message: "`pam_tally2 -r` clears the PAM failed-auth counter — masks ongoing brute force. Log the prior count and alert before resetting.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1592")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
