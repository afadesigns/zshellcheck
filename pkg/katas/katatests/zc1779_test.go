// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1779(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `az role assignment create --role Reader --assignee u --scope s`",
			input:    `az role assignment create --role Reader --assignee u --scope s`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `az role assignment list` (read-only)",
			input:    `az role assignment list`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `az role assignment create --role Owner --assignee u --scope s`",
			input: `az role assignment create --role Owner --assignee u --scope s`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1779",
					Message: "`az role assignment create --role Owner` grants a top-of-chain role. Pick a narrower built-in role (`Reader`, specific-service Contributor) or a custom role whose permission list you can enumerate.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `az role assignment create --role=\"User Access Administrator\" --assignee u --scope s`",
			input: `az role assignment create --role="User Access Administrator" --assignee u --scope s`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1779",
					Message: "`az role assignment create --role User Access Administrator` grants a top-of-chain role. Pick a narrower built-in role (`Reader`, specific-service Contributor) or a custom role whose permission list you can enumerate.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1779")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
