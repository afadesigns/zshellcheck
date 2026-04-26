// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1495(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ulimit -n 10240",
			input:    `ulimit -n 10240`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — ulimit -c 0 (disable)",
			input:    `ulimit -c 0`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ulimit -c unlimited",
			input: `ulimit -c unlimited`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1495",
					Message: "`ulimit -c unlimited` exposes setuid-process memory via core dumps. Leave the distro default and use systemd-coredump if you need post-mortems.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1495")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
