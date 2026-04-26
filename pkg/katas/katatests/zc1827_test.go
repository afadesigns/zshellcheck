// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1827(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `npm deprecate mypkg@1.2.3 'use 1.2.4'`",
			input:    `npm deprecate mypkg@1.2.3 'use 1.2.4'`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `npm publish`",
			input:    `npm publish`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `npm unpublish mypkg@1.2.3`",
			input: `npm unpublish mypkg@1.2.3`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1827",
					Message: "`npm unpublish` removes a published version — every downstream that pinned it fails to install on next CI run (the left-pad pattern). Use `npm deprecate PKG@VERSION 'reason'` so the version stays resolvable with a warning.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1827")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
