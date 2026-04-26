// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1057(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "no ls assignment",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:     "safe array assignment",
			input:    `files=(*)`,
			expected: []katas.Violation{},
		},
		{
			name:     "echo not assignment",
			input:    `echo $(ls)`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1057")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
