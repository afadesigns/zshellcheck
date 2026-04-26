// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1605(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — read-only debugfs",
			input:    `debugfs $DEV`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — debugfs -R command (read-only)",
			input:    `debugfs -R "stat foo.txt" $DEV`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — debugfs -w $DEV",
			input: `debugfs -w $DEV`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1605",
					Message: "`debugfs -w` writes to the filesystem outside the kernel's normal path — journal bypassed, locks ignored. Keep it as an interactive rescue tool, not a script path.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — debugfs -w /dev/loop0",
			input: `debugfs -w /dev/loop0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1605",
					Message: "`debugfs -w` writes to the filesystem outside the kernel's normal path — journal bypassed, locks ignored. Keep it as an interactive rescue tool, not a script path.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1605")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
