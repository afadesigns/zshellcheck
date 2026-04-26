// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1267(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid df -P",
			input:    `df -P /`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid df -h without -P",
			input: `df -h /`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1267",
					Message: "Use `df -P` for script-safe output. `df -h` format varies across systems and may split long device names across lines.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1267")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
