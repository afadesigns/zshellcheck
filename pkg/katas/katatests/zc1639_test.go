package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1639(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — non-credential header",
			input:    `curl -H "Content-Type: application/json" https://api`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — -H @file (read from file)",
			input:    `curl -H @/run/secrets/auth_header https://api`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — Authorization Bearer",
			input: `curl -H "Authorization: Bearer $TOKEN" https://api`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1639",
					Message: "`curl -H \"Authorization: Bearer $TOKEN\"` places the credential in argv — visible via `ps`. Use `-H @FILE` or `--config FILE` with 0600 perms.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — X-Api-Key",
			input: `curl -H "X-Api-Key: $KEY" https://api`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1639",
					Message: "`curl -H \"X-Api-Key: $KEY\"` places the credential in argv — visible via `ps`. Use `-H @FILE` or `--config FILE` with 0600 perms.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1639")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
