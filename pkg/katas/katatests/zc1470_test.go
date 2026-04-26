// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1470(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — git clone https",
			input:    `git clone https://example.com/x.git`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — git config http.sslVerify true",
			input:    `git config http.sslVerify true`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — git config http.sslVerify false",
			input: `git config http.sslVerify false`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1470",
					Message: "`http.sslVerify=false` disables TLS verification — any MITM swaps the clone for attacker code. Fix the CA instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — git config --global http.sslVerify false",
			input: `git config --global http.sslVerify false`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1470",
					Message: "`http.sslVerify=false` disables TLS verification — any MITM swaps the clone for attacker code. Fix the CA instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — git -c http.sslVerify=false clone",
			input: `git -c http.sslVerify=false clone https://example.com/x.git`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1470",
					Message: "`http.sslVerify=false` disables TLS verification — any MITM swaps the clone for attacker code. Fix the CA instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1470")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
