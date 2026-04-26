// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1683(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — https registry",
			input:    `npm config set registry https://registry.npmjs.org/`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — npm config set strict-ssl",
			input:    `npm config set strict-ssl false`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — npm config set registry http://...",
			input: `npm config set registry http://internal.example.com/`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1683",
					Message: "`npm config set registry http://internal.example.com/` uses plaintext HTTP — any proxy / CDN can rewrite tarballs. Use `https://` and a custom CA via `NODE_EXTRA_CA_CERTS` if needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — yarn config set registry http://...",
			input: `yarn config set registry http://internal/`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1683",
					Message: "`yarn config set registry http://internal/` uses plaintext HTTP — any proxy / CDN can rewrite tarballs. Use `https://` and a custom CA via `NODE_EXTRA_CA_CERTS` if needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1683")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
