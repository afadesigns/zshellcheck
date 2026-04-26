// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1661(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — curl with real CA bundle",
			input:    `curl https://example.com --cacert /etc/ssl/certs/ca-certificates.crt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — curl without cacert",
			input:    `curl https://example.com -o out`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — curl URL --cacert /dev/null",
			input: `curl https://example.com --cacert /dev/null`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1661",
					Message: "`curl --cacert /dev/null` feeds curl an empty trust store — most TLS backends then accept any peer cert. Use a real bundle or `--pinnedpubkey`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — curl -s URL --capath /dev/null",
			input: `curl -s https://example.com --capath /dev/null`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1661",
					Message: "`curl --cacert /dev/null` feeds curl an empty trust store — most TLS backends then accept any peer cert. Use a real bundle or `--pinnedpubkey`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1661")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
