package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1712(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — vault kv put with @file",
			input:    `vault kv put secret/app @secret.json`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — vault kv put with stdin sentinel",
			input:    `vault kv put secret/app password=-`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — vault kv put with non-secret key",
			input:    `vault kv put secret/app environment=prod`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — vault kv put password=hunter2",
			input: `vault kv put secret/app password=hunter2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1712",
					Message: "`vault kv password=hunter2` puts the secret value in argv — visible to every local user. Use `password=@FILE` or `password=-` to read from disk / stdin.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — vault write secret/app api_key=ABC",
			input: `vault write secret/app api_key=ABC123`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1712",
					Message: "`vault write api_key=ABC123` puts the secret value in argv — visible to every local user. Use `api_key=@FILE` or `api_key=-` to read from disk / stdin.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1712")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
