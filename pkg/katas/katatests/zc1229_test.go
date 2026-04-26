// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1229(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid rsync",
			input:    `rsync -az src/ user@host:dst/`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid scp",
			input: `scp file.tar.gz user@host:/tmp/`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1229",
					Message: "Prefer `rsync -az` over `scp` for file transfers. `rsync` supports delta transfers, resume, and is more efficient.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1229")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
