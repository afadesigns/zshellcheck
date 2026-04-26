// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1744(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `kubectl port-forward pod/mypod 8080:8080` (loopback default)",
			input:    `kubectl port-forward pod/mypod 8080:8080`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `kubectl port-forward --address 127.0.0.1 pod/mypod 8080:8080`",
			input:    `kubectl port-forward --address 127.0.0.1 pod/mypod 8080:8080`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `kubectl port-forward pod/mypod 8080:8080 --address 0.0.0.0`",
			input: `kubectl port-forward pod/mypod 8080:8080 --address 0.0.0.0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1744",
					Message: "`kubectl port-forward --address 0.0.0.0` binds the local end of the tunnel on every interface — anyone on the LAN / VPN can reach the pod. Drop `--address` (loopback default) or pick a trusted-network interface IP.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `kubectl port-forward pod/mypod 8080:8080 --address=0.0.0.0`",
			input: `kubectl port-forward pod/mypod 8080:8080 --address=0.0.0.0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1744",
					Message: "`kubectl port-forward --address=0.0.0.0` binds the local end of the tunnel on every interface — anyone on the LAN / VPN can reach the pod. Drop `--address` (loopback default) or pick a trusted-network interface IP.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1744")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
