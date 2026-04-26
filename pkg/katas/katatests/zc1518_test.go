// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1518(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — bash -c 'cmd'",
			input:    `bash -c 'true'`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — bash -p -c 'cmd'",
			input: `bash -p -c 'true'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1518",
					Message: "`bash -p` keeps the privileged environment on a setuid wrapper — almost never needed, audit and remove.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1518")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
