// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1888(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `aws iam list-access-keys`",
			input:    `aws iam list-access-keys`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `aws iam get-role --role-name foo`",
			input:    `aws iam get-role --role-name foo`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `aws iam create-access-key --user-name ci-bot`",
			input: `aws iam create-access-key --user-name ci-bot`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1888",
					Message: "`aws iam create-access-key` mints a long-lived `AKIA.../secret` — prefer short-lived creds via instance profiles, IRSA, Lambda roles, or OIDC federation. If static keys are unavoidable, store in Secrets Manager and rotate.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1888")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
