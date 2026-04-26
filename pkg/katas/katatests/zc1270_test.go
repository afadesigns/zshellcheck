// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1270(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid mktemp usage",
			input:    `local tmpfile=$(mktemp)`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid touch with non-tmp path",
			input:    `touch /var/log/myapp.log`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid dynamic tmp path",
			input:    `touch /tmp/$USER-cache`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid hardcoded tmp touch",
			input: `touch /tmp/myfile.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1270",
					Message: "Use `mktemp` instead of hardcoded `/tmp/myfile.txt`. Hardcoded `/tmp` paths are vulnerable to symlink attacks and race conditions.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid hardcoded tmp cat",
			input: `cat /tmp/output.log`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1270",
					Message: "Use `mktemp` instead of hardcoded `/tmp/output.log`. Hardcoded `/tmp` paths are vulnerable to symlink attacks and race conditions.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1270")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
