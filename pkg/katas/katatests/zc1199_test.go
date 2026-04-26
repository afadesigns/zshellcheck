// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1199(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid curl check",
			input:    `curl -s http://localhost:8080`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid telnet",
			input: `telnet localhost 8080`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1199",
					Message: "Avoid `telnet` in scripts — it is interactive and insecure. Use `curl` for HTTP checks or `zmodload zsh/net/tcp` for port testing.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1199")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
