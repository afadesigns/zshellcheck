// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1962(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `kustomize build overlays/prod`",
			input:    `kustomize build overlays/prod`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `kustomize build . --load-restrictor=LoadRestrictionsRootOnly`",
			input:    `kustomize build . --load-restrictor=LoadRestrictionsRootOnly`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `kustomize build . --load-restrictor=LoadRestrictionsNone`",
			input: `kustomize build . --load-restrictor=LoadRestrictionsNone`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1962",
					Message: "`kustomize build --load-restrictor=LoadRestrictionsNone` drops path-root restriction — untrusted overlays can reference `../../secrets/prod.env` and pull them into the render. Keep the default; vendor sibling files into the overlay.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1962")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
