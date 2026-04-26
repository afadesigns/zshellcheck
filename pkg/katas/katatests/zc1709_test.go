// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1709(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — htpasswd -i (read stdin)",
			input:    `htpasswd -i /etc/nginx/.htpasswd user`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — htpasswd interactive (prompts)",
			input:    `htpasswd /etc/nginx/.htpasswd user`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — htpasswd -b user secret",
			input: `htpasswd -b /etc/nginx/.htpasswd user secret`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1709",
					Message: "`htpasswd -b USER PASSWORD` puts the password in argv — visible via `ps` / `/proc/PID/cmdline`. Use `htpasswd -i FILE USER` with the password piped on stdin instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — htpasswd -bB combined flags",
			input: `htpasswd -bB /etc/nginx/.htpasswd user secret`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1709",
					Message: "`htpasswd -b USER PASSWORD` puts the password in argv — visible via `ps` / `/proc/PID/cmdline`. Use `htpasswd -i FILE USER` with the password piped on stdin instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1709")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
