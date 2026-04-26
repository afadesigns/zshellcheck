// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1691(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — rsync without --remove-source-files",
			input:    `rsync -av src/ dst/`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — rsync --delete (different flag)",
			input:    `rsync -av --delete src/ dst/`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — rsync --remove-source-files (local)",
			input: `rsync -av --remove-source-files src/ dst/`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1691",
					Message: "`rsync --remove-source-files` deletes SRC on optimistic per-file success — verify DST after the transfer and `rm` explicitly instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — rsync --remove-source-files (remote)",
			input: `rsync -a --remove-source-files host:src dst`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1691",
					Message: "`rsync --remove-source-files` deletes SRC on optimistic per-file success — verify DST after the transfer and `rm` explicitly instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1691")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
