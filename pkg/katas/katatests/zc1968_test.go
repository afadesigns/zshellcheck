// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1968(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `dnf versionlock list`",
			input:    `dnf versionlock list`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `dnf update --exclude=kernel`",
			input:    `dnf update --exclude=kernel`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `dnf versionlock add kernel`",
			input: `dnf versionlock add kernel`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1968",
					Message: "`dnf versionlock add` pins the rpm — blocks future CVE fixes for glibc/openssl/kernel. Prefer `--exclude` on a single transaction and schedule a `versionlock delete` review.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `yum versionlock add openssl`",
			input: `yum versionlock add openssl`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1968",
					Message: "`yum versionlock add` pins the rpm — blocks future CVE fixes for glibc/openssl/kernel. Prefer `--exclude` on a single transaction and schedule a `versionlock delete` review.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1968")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
