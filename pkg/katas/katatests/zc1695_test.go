// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1695(t *testing.T) {
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
			name:     "valid — terraform state mv (tracked rename)",
			input:    `terraform state mv old new`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — terraform state rm",
			input: `terraform state rm module.app.aws_instance.x`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1695",
					Message: "`terraform state rm` mutates shared state outside plan/apply — use `terraform import` or `apply -replace=ADDR` instead, and review / back up first.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — tofu state push",
			input: `tofu state push local.tfstate`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1695",
					Message: "`tofu state push` mutates shared state outside plan/apply — use `terraform import` or `apply -replace=ADDR` instead, and review / back up first.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1695")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
