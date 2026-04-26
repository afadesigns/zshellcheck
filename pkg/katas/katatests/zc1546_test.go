// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1546(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — kubectl delete pod foo",
			input:    `kubectl delete pod foo`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — kubectl delete pod foo --force (grace-period not 0)",
			input:    `kubectl delete pod foo --force`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — kubectl delete pod foo --force --grace-period=0",
			input: `kubectl delete pod foo --force --grace-period=0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1546",
					Message: "`kubectl delete --force --grace-period=0` skips PreStop hooks and kubelet drain — corrupts StatefulSet state. Use standard delete.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — oc delete pod foo --force --grace-period=0",
			input: `oc delete pod foo --force --grace-period=0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1546",
					Message: "`kubectl delete --force --grace-period=0` skips PreStop hooks and kubelet drain — corrupts StatefulSet state. Use standard delete.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1546")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
