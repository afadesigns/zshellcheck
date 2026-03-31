package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1228(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid ssh with BatchMode",
			input:    `ssh -o BatchMode=yes user@host`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid ssh without policy",
			input: `ssh user@host ls`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1228",
					Message: "Use `ssh -o BatchMode=yes` or `-o StrictHostKeyChecking=accept-new` in scripts. Without these, ssh may prompt interactively and hang.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1228")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
