// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1058(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "sudo without redirection",
			input:    `sudo apt install vim`,
			expected: []katas.Violation{},
		},
		{
			name:     "not sudo command",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1058")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
