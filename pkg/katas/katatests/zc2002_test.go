// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC2002(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `crictl ps`",
			input:    `crictl ps`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `crictl rmi $IMAGE_ID`",
			input:    `crictl rmi $IMAGE_ID`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `crictl rmi -a`",
			input: `crictl rmi -a`,
			expected: []katas.Violation{
				{
					KataID:  "ZC2002",
					Message: "`crictl rmi -a` talks to the node CRI directly, under the kubelet — images/containers backing running pods disappear, kubelet must re-pull or re-run. Route through `kubectl drain`/`delete pod`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `crictl rm -af`",
			input: `crictl rm -af`,
			expected: []katas.Violation{
				{
					KataID:  "ZC2002",
					Message: "`crictl rm -af` talks to the node CRI directly, under the kubelet — images/containers backing running pods disappear, kubelet must re-pull or re-run. Route through `kubectl drain`/`delete pod`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC2002")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
