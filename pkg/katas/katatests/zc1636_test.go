// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1636(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — virsh shutdown",
			input:    `virsh shutdown my-vm`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — virsh destroy --graceful",
			input:    `virsh destroy --graceful my-vm`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — virsh destroy my-vm",
			input: `virsh destroy my-vm`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1636",
					Message: "`virsh destroy` yanks power from the VM — filesystem corruption risk. Use `virsh shutdown` for graceful stop, or `virsh destroy --graceful` as a timed fallback.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — virsh destroy $DOM",
			input: `virsh destroy $DOM`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1636",
					Message: "`virsh destroy` yanks power from the VM — filesystem corruption risk. Use `virsh shutdown` for graceful stop, or `virsh destroy --graceful` as a timed fallback.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1636")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
