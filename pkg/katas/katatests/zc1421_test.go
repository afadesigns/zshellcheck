package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1421(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — chpasswd -e (encrypted)",
			input:    `chpasswd -e`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — chpasswd -c (plaintext)",
			input: `chpasswd -c SHA512`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1421",
					Message: "`chpasswd` without `-e`/`--encrypted` accepts plaintext passwords — avoid piping cleartext credentials into the process tree. Use a password hash (`-e`) or a credentials store.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1421")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
