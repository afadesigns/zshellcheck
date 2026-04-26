// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC2000(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `kubectl taint nodes $NODE key=value:NoSchedule` (gentle)",
			input:    `kubectl taint nodes $NODE key=value:NoSchedule`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `kubectl drain $NODE`",
			input:    `kubectl drain $NODE`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `kubectl taint nodes $NODE key=value:NoExecute`",
			input: `kubectl taint nodes $NODE key=value:NoExecute`,
			expected: []katas.Violation{
				{
					KataID:  "ZC2000",
					Message: "`kubectl taint nodes … :NoExecute` evicts every non-tolerating pod immediately — a typo on `--all` nodes empties the cluster. Prefer `kubectl drain $NODE` or a `:NoSchedule` taint for gentle drain.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC2000")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
