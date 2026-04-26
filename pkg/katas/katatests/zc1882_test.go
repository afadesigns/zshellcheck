// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1882(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `sudo /usr/local/bin/setup.sh`",
			input:    `sudo /usr/local/bin/setup.sh`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `sudo -i /usr/local/bin/setup.sh`",
			input:    `sudo -i /usr/local/bin/setup.sh`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `sudo su -c \"cmd\"`",
			input:    `sudo su -c "cmd"`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `sudo -s`",
			input: `sudo -s`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1882",
					Message: "`sudo -s` spawns an interactive root shell — in a script either hangs on stdin or drains the rest of the file into root's shell. Pass the command to sudo: `sudo /path/to/cmd arg …`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `sudo su -`",
			input: `sudo su -`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1882",
					Message: "`sudo su` spawns an interactive root shell — in a script either hangs on stdin or drains the rest of the file into root's shell. Pass the command to sudo: `sudo /path/to/cmd arg …`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `sudo bash`",
			input: `sudo bash`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1882",
					Message: "`sudo bash` spawns an interactive root shell — in a script either hangs on stdin or drains the rest of the file into root's shell. Pass the command to sudo: `sudo /path/to/cmd arg …`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1882")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
