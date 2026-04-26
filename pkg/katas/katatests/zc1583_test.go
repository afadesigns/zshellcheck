// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1583(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — find without -delete",
			input:    `find /tmp -name '*.log'`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — find -maxdepth 2 -delete",
			input:    `find /tmp -maxdepth 2 -name '*.log' -delete`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — find -xdev -delete",
			input:    `find /var -xdev -name '*.log' -delete`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — find /tmp -name '*.log' -delete",
			input: `find /tmp -name '*.log' -delete`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1583",
					Message: "`find -delete` without `-maxdepth` / `-xdev` / `-prune` walks the whole tree. Scope the depth (e.g. `-maxdepth 2`) and dry-run first.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1583")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
