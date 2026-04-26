// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1814(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `dpkg -i pkg.deb` (no force)",
			input:    `dpkg -i pkg.deb`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `dpkg --force-overwrite -i pkg.deb` (specific force)",
			input:    `dpkg --force-overwrite -i pkg.deb`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `dpkg -i --force-all pkg.deb`",
			input: `dpkg -i --force-all pkg.deb`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1814",
					Message: "`dpkg --force-all` enables every `--force-*` option at once — overwrite, unsigned, downgrade, essential-removal, broken-deps. Drop it and spell out only the specific `--force-<option>` you need, or fix the upstream conflict.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `apt-get -o Dpkg::Options::=--force-all install pkg`",
			input: `apt-get -o Dpkg::Options::=--force-all install pkg`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1814",
					Message: "`dpkg --force-all` enables every `--force-*` option at once — overwrite, unsigned, downgrade, essential-removal, broken-deps. Drop it and spell out only the specific `--force-<option>` you need, or fix the upstream conflict.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1814")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
