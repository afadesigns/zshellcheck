// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1434(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — swapoff specific file",
			input:    `swapoff swap.img`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — swapoff -a",
			input: `swapoff -a`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1434",
					Message: "`swapoff -a` disables ALL swap areas — risks OOM on memory-constrained hosts. Disable specific swaps (`swapoff /swapfile`) after checking `free -m`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1434")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
