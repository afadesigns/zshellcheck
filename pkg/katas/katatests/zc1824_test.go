// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1824(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `kubectl drain node-1 --ignore-daemonsets`",
			input:    `kubectl drain node-1 --ignore-daemonsets`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `kubectl cordon node-1`",
			input:    `kubectl cordon node-1`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `kubectl drain node-1 --disable-eviction`",
			input: `kubectl drain node-1 --disable-eviction`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1824",
					Message: "`kubectl drain --disable-eviction` deletes pods via raw API DELETE — PodDisruptionBudgets are ignored and the workload owner's availability contract is voided. Fix the blocking PDB instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1824")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
