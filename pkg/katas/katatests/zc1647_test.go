// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1647(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — local file",
			input:    `kubectl apply -f ./manifest.yaml`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — stdin",
			input:    `kubectl apply -f -`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — kubectl apply -f https://example.com/m.yaml",
			input: `kubectl apply -f https://example.com/manifest.yaml`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1647",
					Message: "`kubectl apply -f https://example.com/manifest.yaml` applies a remote manifest — verify digest first. Download, check SHA256, then apply the local file.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — kubectl create -f http://insecure/m.yaml",
			input: `kubectl create -f http://insecure/m.yaml`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1647",
					Message: "`kubectl create -f http://insecure/m.yaml` applies a remote manifest — verify digest first. Download, check SHA256, then apply the local file.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1647")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
