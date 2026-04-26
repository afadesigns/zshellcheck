// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1793(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `kubectl certificate deny CSR_NAME`",
			input:    `kubectl certificate deny CSR_NAME`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `kubectl get csr`",
			input:    `kubectl get csr`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `kubectl certificate approve CSR_NAME`",
			input: `kubectl certificate approve CSR_NAME`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1793",
					Message: "`kubectl certificate approve` signs the identity embedded in the CSR — a `system:masters` request becomes cluster admin. Decode with `openssl req -text` first; use `kubectl certificate deny` otherwise.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1793")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
