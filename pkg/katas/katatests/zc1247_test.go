// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1247(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid chmod 755",
			input:    `chmod 755 script.sh`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid chmod 2755",
			input: `chmod 2755 /usr/local/bin/tool`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1247",
					Message: "Avoid `chmod +s` — setuid/setgid bits create privilege escalation risks. Use `sudo`, capabilities (`setcap`), or dedicated service accounts instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid chmod 4755",
			input: `chmod 4755 /usr/local/bin/tool`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1247",
					Message: "Avoid `chmod +s` — setuid/setgid bits create privilege escalation risks. Use `sudo`, capabilities (`setcap`), or dedicated service accounts instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1247")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
