// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1759(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `vault login -` (reads token from stdin)",
			input:    `vault login -`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `vault login -method=userpass username=alice` (no secret key)",
			input:    `vault login -method=userpass username=alice`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `vault login mytoken` (positional token)",
			input: `vault login mytoken`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1759",
					Message: "`vault login mytoken` puts the Vault credential in argv — visible in `ps`, `/proc`, history, Vault audit log. Use `vault login -` with stdin or source `VAULT_TOKEN` from a secrets file.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `vault login -method=userpass username=alice password=hunter2`",
			input: `vault login -method=userpass username=alice password=hunter2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1759",
					Message: "`vault login password=hunter2` puts the Vault credential in argv — visible in `ps`, `/proc`, history, Vault audit log. Use `vault login -` with stdin or source `VAULT_TOKEN` from a secrets file.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `vault auth mytoken`",
			input: `vault auth mytoken`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1759",
					Message: "`vault auth mytoken` puts the Vault credential in argv — visible in `ps`, `/proc`, history, Vault audit log. Use `vault login -` with stdin or source `VAULT_TOKEN` from a secrets file.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1759")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
