// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1551(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — helm install chart",
			input:    `helm install foo bitnami/nginx`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — helm install chart --skip-crds",
			input: `helm install foo bitnami/nginx --skip-crds`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1551",
					Message: "`helm --skip-crds` installs .Release objects without their CRDs — custom resources fail validation. Install CRDs first.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — helm upgrade chart --skip-crds",
			input: `helm upgrade foo bitnami/nginx --skip-crds`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1551",
					Message: "`helm --skip-crds` installs .Release objects without their CRDs — custom resources fail validation. Install CRDs first.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1551")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
