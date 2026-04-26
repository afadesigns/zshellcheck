// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1576(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — terraform apply",
			input:    `terraform apply`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — terraform apply -target=module.foo",
			input: `terraform apply -target=module.foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1576",
					Message: "`terraform -target=module.foo` bypasses dependency order — documented as incident response tool only. Re-run without -target or split root modules.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — terraform apply -target module.foo",
			input: `terraform apply -target module.foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1576",
					Message: "`terraform -target module.foo` bypasses dependency order — documented as incident response tool only. Re-run without -target or split root modules.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — tofu destroy -target=aws_instance.web",
			input: `tofu destroy -target=aws_instance.web`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1576",
					Message: "`terraform -target=aws_instance.web` bypasses dependency order — documented as incident response tool only. Re-run without -target or split root modules.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1576")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
