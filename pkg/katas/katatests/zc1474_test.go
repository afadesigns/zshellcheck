package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1474(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ssh-keygen -N passphrase",
			input:    `ssh-keygen -N secretpass -f key`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — ssh-keygen without -N",
			input:    `ssh-keygen -t ed25519 -f key`,
			expected: []katas.Violation{},
		},
		{
			name:  `invalid — ssh-keygen -N "" -f key`,
			input: `ssh-keygen -N "" -f key`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1474",
					Message: "`ssh-keygen -N \"\"` generates a passwordless key — anything that reads the file can use it. Use a passphrase or ssh-agent/HSM.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  `invalid — ssh-keygen -N '' -f key`,
			input: `ssh-keygen -N '' -f key`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1474",
					Message: "`ssh-keygen -N \"\"` generates a passwordless key — anything that reads the file can use it. Use a passphrase or ssh-agent/HSM.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1474")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
