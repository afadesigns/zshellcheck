// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1040(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "safe glob with nullglob qualifier",
			input:    `for f in *.txt(N); do echo $f; done`,
			expected: []katas.Violation{},
		},
		{
			name:     "no glob pattern",
			input:    `for f in a b c; do echo $f; done`,
			expected: []katas.Violation{},
		},
		{
			name:     "arithmetic for loop",
			input:    `for (( i=0; i<10; i++ )); do echo $i; done`,
			expected: []katas.Violation{},
		},
		{
			name:     "quoted string is not a glob",
			input:    `for f in "*.txt"; do echo $f; done`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1040")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
