package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1432(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — chattr -i (remove)",
			input:    `chattr -i file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — chattr +i",
			input: `chattr +i secret.conf`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1432",
					Message: "`chattr +i`/`+a` sets immutable/append-only — blocks later cleanup. Document the `-i`/`-a` cleanup path or reconsider if really needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — chattr +a",
			input: `chattr +a log.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1432",
					Message: "`chattr +i`/`+a` sets immutable/append-only — blocks later cleanup. Document the `-i`/`-a` cleanup path or reconsider if really needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1432")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
