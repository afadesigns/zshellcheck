// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1762(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `kubeadm join ... --discovery-token-ca-cert-hash sha256:xxx`",
			input:    `kubeadm join 10.0.0.1:6443 --token abc --discovery-token-ca-cert-hash sha256:xxx`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `kubeadm token create`",
			input:    `kubeadm token create --print-join-command`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `kubeadm join ... --discovery-token-unsafe-skip-ca-verification`",
			input: `kubeadm join 10.0.0.1:6443 --token abc --discovery-token-unsafe-skip-ca-verification`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1762",
					Message: "`kubeadm join --discovery-token-unsafe-skip-ca-verification` skips CA verification of the control-plane — MITM steals the bootstrap token. Pin the CA with `--discovery-token-ca-cert-hash sha256:<digest>`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1762")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
