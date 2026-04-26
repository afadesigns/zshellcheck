// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1429(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — umount /mnt",
			input:    `umount /mnt/disk`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — umount -f",
			input: `umount -f /mnt/disk`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1429",
					Message: "`umount -f`/`-l` force/lazy unmount masks the underlying 'busy' error. Find open files with `lsof` / `fuser -m` and close them properly.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — umount -l",
			input: `umount -l /mnt/disk`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1429",
					Message: "`umount -f`/`-l` force/lazy unmount masks the underlying 'busy' error. Find open files with `lsof` / `fuser -m` and close them properly.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1429")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
