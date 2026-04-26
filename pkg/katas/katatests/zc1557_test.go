// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1557(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — kubeadm init",
			input:    `kubeadm init`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — kubeadm reset (no -f)",
			input:    `kubeadm reset`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — kubeadm reset -f",
			input: `kubeadm reset -f`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1557",
					Message: "`kubeadm reset -f` skips the confirmation and wipes /etc/kubernetes / kubelet state. Drain and remove the node first.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — kubeadm reset --force",
			input: `kubeadm reset --force`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1557",
					Message: "`kubeadm reset -f` skips the confirmation and wipes /etc/kubernetes / kubelet state. Drain and remove the node first.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1557")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
