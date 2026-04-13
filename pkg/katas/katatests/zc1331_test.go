package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1331(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid match usage",
			input:    `echo $match`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid BASH_REMATCH usage",
			input: `echo $BASH_REMATCH`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1331",
					Message: "Avoid `$BASH_REMATCH` in Zsh — use `$match` array and `$MATCH` for regex captures instead.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1331")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
