// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1480(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — terraform plan",
			input:    `terraform plan`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — terraform apply (interactive)",
			input:    `terraform apply`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — terraform apply -auto-approve",
			input: `terraform apply -auto-approve`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1480",
					Message: "`terraform apply -auto-approve` skips plan review. Gate behind a branch/env check or require manual approval.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — terraform destroy --auto-approve",
			input: `terraform destroy --auto-approve`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1480",
					Message: "`terraform destroy --auto-approve` skips plan review. Gate behind a branch/env check or require manual approval.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — tofu apply -auto-approve",
			input: `tofu apply -auto-approve`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1480",
					Message: "`tofu apply -auto-approve` skips plan review. Gate behind a branch/env check or require manual approval.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1480")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
