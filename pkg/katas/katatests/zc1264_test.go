// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1264(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid dnf",
			input:    `dnf install -y curl`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid yum",
			input: `yum install -y curl`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1264",
					Message: "Use `dnf` instead of `yum`. `yum` is deprecated on modern Fedora and RHEL; `dnf` has better dependency resolution.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1264")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
