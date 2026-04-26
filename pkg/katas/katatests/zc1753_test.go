// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1753(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `rclone delete myremote:bucket/path`",
			input:    `rclone delete myremote:bucket/path`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `rclone sync src dst`",
			input:    `rclone sync src dst`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `rclone purge myremote:bucket/path`",
			input: `rclone purge myremote:bucket/path`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1753",
					Message: "`rclone purge` removes every object under the remote path with no dry-run or soft-delete. Preview with `rclone lsf` / `rclone delete --dry-run` and prefer narrower `rclone delete`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1753")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
