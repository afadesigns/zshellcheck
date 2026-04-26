// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1662(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — pkexec direct command",
			input:    `pkexec /usr/bin/systemctl restart unit`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — pkexec apt install",
			input:    `pkexec /usr/bin/apt install foo`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — pkexec env DISPLAY=... cmd",
			input: `pkexec env DISPLAY=:0 /usr/bin/cmd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1662",
					Message: "`pkexec env VAR=VAL CMD` hands the root session a caller-controlled environment — use a polkit rule or `systemd-run --user` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — pkexec env PATH=/tmp cmd",
			input: `pkexec env PATH=/tmp /usr/bin/cmd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1662",
					Message: "`pkexec env VAR=VAL CMD` hands the root session a caller-controlled environment — use a polkit rule or `systemd-run --user` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1662")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
