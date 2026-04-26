// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1507(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — rsync without archive / -l",
			input:    `rsync -rv src/ dst/`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — rsync -a --safe-links src/ dst/",
			input:    `rsync -a --safe-links src/ dst/`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — rsync --no-links src/ dst/",
			input:    `rsync -a --no-links src/ dst/`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — rsync -a src/ dst/",
			input: `rsync -a src/ dst/`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1507",
					Message: "`rsync` preserving symlinks without `--safe-links` follows ones pointing outside the source tree — path traversal vector. Add `--safe-links` or `--copy-unsafe-links`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — rsync -al src/ dst/",
			input: `rsync -al src/ dst/`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1507",
					Message: "`rsync` preserving symlinks without `--safe-links` follows ones pointing outside the source tree — path traversal vector. Add `--safe-links` or `--copy-unsafe-links`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1507")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
