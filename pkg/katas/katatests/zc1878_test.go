// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1878(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `kubectl apply -f manifest.yaml`",
			input:    `kubectl apply -f manifest.yaml`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `kubectl apply --server-side -f manifest.yaml`",
			input:    `kubectl apply --server-side -f manifest.yaml`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `kubectl apply --server-side --force-conflicts -f manifest.yaml`",
			input: `kubectl apply --server-side --force-conflicts -f manifest.yaml`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1878",
					Message: "`kubectl apply --force-conflicts` grabs ownership of every conflicting field from other controllers (HPA, cert-manager, sidecar injectors). Resolve the conflict instead — drop the disputed fields or hand off via managed-field edit.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1878")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
