// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1891(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `kubectl config view`",
			input:    `kubectl config view`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `kubectl config view -o jsonpath='{.current-context}'`",
			input:    `kubectl config view -o jsonpath='{.current-context}'`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `kubectl config view --raw`",
			input: `kubectl config view --raw`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1891",
					Message: "`kubectl config view --raw` prints the full kubeconfig including client-certificate/key-data and bearer tokens — any script-captured stdout exfiltrates the creds. Emit the specific field with `-o jsonpath='…'`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `kubectl config view -R`",
			input: `kubectl config view -R`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1891",
					Message: "`kubectl config view --raw` prints the full kubeconfig including client-certificate/key-data and bearer tokens — any script-captured stdout exfiltrates the creds. Emit the specific field with `-o jsonpath='…'`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1891")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
