// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1393(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — echo $RANDOM",
			input:    `echo $RANDOM`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $SRANDOM",
			input: `echo $SRANDOM`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1393",
					Message: "`$SRANDOM` is Bash 5.1+. In Zsh read `/dev/urandom` directly or use an external (`openssl rand`) for secure random integers.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1393")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
