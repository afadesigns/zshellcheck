// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1567(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — python -m http.server --bind 127.0.0.1",
			input:    `python -m http.server --bind 127.0.0.1 8080`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — python -m http.server -b 127.0.0.1",
			input:    `python -m http.server -b 127.0.0.1`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — python -m venv myenv",
			input:    `python -m venv myenv`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — python -m http.server",
			input: `python -m http.server`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1567",
					Message: "`python -m http.server` without `--bind` defaults to 0.0.0.0 — exposes the cwd to every network the host sees. Add `--bind 127.0.0.1`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — python3 -m http.server 8080",
			input: `python3 -m http.server 8080`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1567",
					Message: "`python -m http.server` without `--bind` defaults to 0.0.0.0 — exposes the cwd to every network the host sees. Add `--bind 127.0.0.1`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — python2 -m SimpleHTTPServer",
			input: `python2 -m SimpleHTTPServer`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1567",
					Message: "`python -m http.server` without `--bind` defaults to 0.0.0.0 — exposes the cwd to every network the host sees. Add `--bind 127.0.0.1`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1567")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
