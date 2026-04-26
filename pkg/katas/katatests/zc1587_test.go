// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1587(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — modprobe nvme",
			input:    `modprobe nvme`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — lsmod",
			input:    `lsmod`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — modprobe -r nvme",
			input: `modprobe -r nvme`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1587",
					Message: "`modprobe -r` unloads an in-use module — the backing subsystem goes offline. Use `systemctl stop` if you meant to stop a service.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — rmmod nvidia",
			input: `rmmod nvidia`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1587",
					Message: "`rmmod` unloads a kernel module — the backing subsystem goes offline. Use `systemctl stop` if you meant to stop a service.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1587")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
