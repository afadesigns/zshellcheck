// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1468(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — apt-get install curl",
			input:    `apt-get install curl`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — apt-get install -y curl",
			input:    `apt-get install -y curl`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — apt-get install --allow-unauthenticated curl",
			input: `apt-get install --allow-unauthenticated curl`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1468",
					Message: "APT installing unsigned or override-policy packages (--allow-unauthenticated) — disables signature verification, MITM-to-root trivial.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — apt-get install --force-yes foo",
			input: `apt-get install --force-yes foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1468",
					Message: "APT installing unsigned or override-policy packages (--force-yes) — disables signature verification, MITM-to-root trivial.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — apt install --allow-downgrades foo",
			input: `apt install --allow-downgrades foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1468",
					Message: "APT installing unsigned or override-policy packages (--allow-downgrades) — disables signature verification, MITM-to-root trivial.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1468")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
