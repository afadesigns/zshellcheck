// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1563(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — swapon -a",
			input:    `swapon -a`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — swapoff specific file",
			input:    `swapoff /swapfile`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — swapoff -a",
			input: `swapoff -a`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1563",
					Message: "`swapoff -a` disables all swap devices — next memory-hungry process hits OOM. Document the trade-off if kubelet requires it.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1563")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
