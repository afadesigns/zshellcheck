// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1740(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `gh release upload v1.0 file.tar.gz`",
			input:    `gh release upload v1.0 file.tar.gz`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `gh release create v1.0 file.tar.gz`",
			input:    `gh release create v1.0 file.tar.gz`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `gh release upload v1.0 file.tar.gz --clobber`",
			input: `gh release upload v1.0 file.tar.gz --clobber`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1740",
					Message: "`gh release upload --clobber` silently replaces an existing asset — a re-run can downgrade the user-facing download. Drop `--clobber` or version the asset name so each upload has a unique slot.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1740")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
