// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1561(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — systemctl isolate multi-user.target",
			input:    `systemctl isolate multi-user.target`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — systemctl start nginx.service",
			input:    `systemctl start nginx.service`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — systemctl isolate rescue.target",
			input: `systemctl isolate rescue.target`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1561",
					Message: "`systemctl isolate rescue.target` kills SSH and most services — console-only recovery. Do not run from a script.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — systemctl isolate emergency.target",
			input: `systemctl isolate emergency.target`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1561",
					Message: "`systemctl isolate emergency.target` kills SSH and most services — console-only recovery. Do not run from a script.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1561")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
