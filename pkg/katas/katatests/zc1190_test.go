// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1190(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid grep -v -e combined",
			input:    `grep -v -e foo -e bar file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid chained grep -v",
			input: `grep -v foo file | grep -v bar`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1190",
					Message: "Combine `grep -v p1 | grep -v p2` into `grep -v -e p1 -e p2`. A single invocation avoids an unnecessary pipeline.",
					Line:    1,
					Column:  18,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1190")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
