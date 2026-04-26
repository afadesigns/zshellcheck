// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1552(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — openssl genrsa 2048",
			input:    `openssl genrsa 2048`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — openssl dhparam 4096",
			input:    `openssl dhparam 4096`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — openssl x509 (not key-producing)",
			input:    `openssl x509 -in cert.pem -noout`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — openssl genrsa 1024",
			input: `openssl genrsa 1024`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1552",
					Message: "`openssl genrsa 1024` uses a weak key/param size — modern baselines require 2048+. Use 2048 or 3072/4096 for long-lived keys.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — openssl dhparam 512",
			input: `openssl dhparam 512`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1552",
					Message: "`openssl dhparam 512` uses a weak key/param size — modern baselines require 2048+. Use 2048 or 3072/4096 for long-lived keys.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1552")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
